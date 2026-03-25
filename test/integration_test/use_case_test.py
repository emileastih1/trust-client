import collections
import time
import requests
import threading
import datetime
import pprint
import argparse

from helper import chain_eth, generate_address, get_random_valid_vasp_domain, \
    get_vasp_uuid_by_domain_name

pp = pprint.PrettyPrinter(indent=4)

requesting_vasp = "testvasp1.com"

PATH_TO_MY_CERT = "test/certs_dev/client_int_1.crt"
PATH_TO_MY_KEY_FILE = "test/certs_dev/client_int_1.key"
PATH_TO_CA_BUNDLE_OR_FALSE = "test/certs_dev/bundle.crt"


def stage_to_test():
    return "local"


def send_request(method, uri, headers, body=None):
    if stage_to_test() == "development":
        base_url = "https://development-bbapi-nlb-travel-rul-503376089e6e6323.elb.us-east-1.amazonaws.com:443"
        if method == "GET":
            return requests.get(url=base_url + uri, headers=headers,
                                cert=(PATH_TO_MY_CERT, PATH_TO_MY_KEY_FILE),
                                verify=PATH_TO_CA_BUNDLE_OR_FALSE)
        elif method == "PUT":
            return requests.put(url=base_url + uri, headers=headers, json=body,
                                cert=(PATH_TO_MY_CERT, PATH_TO_MY_KEY_FILE),
                                verify=PATH_TO_CA_BUNDLE_OR_FALSE)
        elif method == "DELETE":
            return requests.delete(url=base_url + uri, headers=headers,
                                   cert=(PATH_TO_MY_CERT, PATH_TO_MY_KEY_FILE),
                                   verify=PATH_TO_CA_BUNDLE_OR_FALSE)
    elif stage_to_test() == "local":
        base_url = "http://localhost:7000"
        if method == "GET":
            return requests.get(url=base_url + uri, headers=headers)
        elif method == "PUT":
            return requests.put(url=base_url + uri, headers=headers, json=body)
        elif method == "DELETE":
            return requests.delete(url=base_url + uri, headers=headers)
    else:
        raise Exception("invalid stage")


def create_address_ownership(req_vasp_domain, hash_addr):
    headers = {
        'content-type': 'application/json',
        'grpc-metadata-x-vasp-dns-name': req_vasp_domain,
    }

    request_uri = "/v1/addresses/{0}".format(hash_addr)
    r = send_request("PUT", uri=request_uri, headers=headers,
                     body={"chain": chain_eth})

    if r.status_code == 200:
        res = r.json()
        registration_id = res["registrationId"]
        return registration_id
    elif r.status_code == 429:
        print("rate limit exception in create_address_ownership")
    else:
        print("create_address_ownership_request", request_uri, headers)
        raise Exception(r.json())


def get_address_ownership(req_vasp_domain, address):  # return vasp_uuid
    headers = {
        'content-type': 'application/json',
        'grpc-metadata-x-vasp-dns-name': req_vasp_domain,
    }
    request_uri = "/v1/addresses/{0}?chain={1}".format(address, chain_eth)
    r = send_request("GET", uri=request_uri, headers=headers)
    if r.status_code == 200:
        res = r.json()
        return res
    elif r.status_code == 429:
        print("rate limit exception in get_address_ownership")
    else:
        print("get_address_ownership_request", address)
        raise Exception(r.json())


def create_address_ownership_proof(req_vasp_domain, hash_addr, iou, proof):
    headers = {
        'content-type': 'application/json',
        'grpc-metadata-x-vasp-dns-name': req_vasp_domain,
    }

    request_uri = "/v1/address_ownership_proofs/{0}".format(hash_addr)
    r = send_request("PUT", uri=request_uri, headers=headers, body={
                            "iou": True,
                            "registration_id": proof,
                            "chain": chain_eth})

    if r.status_code == 200:
        res = r.json()
        registration_id = res["registrationId"]
        return registration_id
    elif r.status_code == 429:
        print("rate limit exception in create_address_ownership_proof")
    else:
        print("create_address_ownership_proof_request", request_uri)
        raise Exception(r.json())


def delete_address_ownership(req_vasp_domain, hash_addr):
    headers = {
        'content-type': 'application/json',
        'grpc-metadata-x-vasp-dns-name': req_vasp_domain,
    }

    request_uri = "/v1/addresses/{0}?chain={1}".format(hash_addr, chain_eth)
    r = send_request("DELETE", uri=request_uri, headers=headers)

    if r.status_code == 200:
        return
    elif r.status_code == 429:
        print("rate limit exception in delete_address_ownership")
    else:
        print("delete_address_ownership_request", request_uri)
        raise Exception(r.json())


########################################################
# workflows
########################################################
total_time_difference = 0.1
print_threshold = 50


def pretty_print(index, request_monitoring):
    if index % print_threshold == 0:
        print("create_address_ownership")
        pp.pprint(dict(request_monitoring["create_address_ownership"]))
        print("create_address_ownership_proof")
        pp.pprint(dict(request_monitoring["create_address_ownership_proof"]))


def request_workflow(thread_index, total_number_of_iterations):
    request_mapping = collections.defaultdict()
    request_monitoring = collections.defaultdict()

    request_monitoring["create_address_ownership"] = collections.defaultdict(
        int)
    request_monitoring["create_address_ownership_proof"] = collections.defaultdict(
        int)

    for i in range(total_number_of_iterations):
        start_time = datetime.datetime.now()
        if i % print_threshold == 0:
            print("thread\t", thread_index, "\tcase creation count\t", i)
        req_vasp_domain = get_random_valid_vasp_domain()
        p_key, addr, req_hash_addr = generate_address()

        # 1. create address ownership
        resp_registration_id = create_address_ownership(
            req_vasp_domain, req_hash_addr)

        request_mapping[req_hash_addr] = {
            "req_vasp_domain": req_vasp_domain,
            "req_vasp_uuid": get_vasp_uuid_by_domain_name(req_vasp_domain),
            "registration_id": resp_registration_id,
            "iou": True
        }
        finish_create_time = datetime.datetime.now()
        diff1 = finish_create_time - start_time
        if diff1.total_seconds() < total_time_difference:
            time.sleep(total_time_difference - diff1.total_seconds())
        time_spent = diff1.total_seconds() * 1000 // 10 * 10
        request_monitoring["create_address_ownership"][str(time_spent)] += 1
        finish_creat_time_start = datetime.datetime.now()

        # 2. submit IOU
        # print("create_address_ownership_proof")
        create_address_ownership_proof(
            req_vasp_domain, req_hash_addr, "true", resp_registration_id)
        finish_create_proof_time = datetime.datetime.now()
        diff2 = finish_create_proof_time - finish_creat_time_start
        if diff2.total_seconds() < total_time_difference:
            time.sleep(total_time_difference - diff2.total_seconds())
        time_spent2 = diff2.total_seconds() * 1000 // 10 * 10
        request_monitoring["create_address_ownership_proof"][str(
            time_spent2)] += 1
    pretty_print(0, request_monitoring)

    # 3. Submit proof
    return request_mapping


def response_verification(thread_index, response_mapping):
    count = 0
    for req_hash_addr, v in response_mapping.items():
        start_time = datetime.datetime.now()
        # 4. get address ownership
        if count % print_threshold == 0:
            print("thread\t", thread_index, "\tverification count\t", count)
        try:
            # print("get_address_ownership")
            response = get_address_ownership(requesting_vasp, req_hash_addr)
            resp_owner_vasp_uuid = response["vasp"]["id"]
            resp_iou = response["addressOwnershipProof"]["iou"]
            finish_get_time = datetime.datetime.now()
            diff1 = finish_get_time - start_time
            if diff1.total_seconds() < total_time_difference:
                time.sleep(total_time_difference - diff1.total_seconds())

            if stage_to_test() == "development":
                assert resp_owner_vasp_uuid == "3c348940-02a2-4198-9a5f-9d127344d46f"
            else:
                assert resp_owner_vasp_uuid == v["req_vasp_uuid"]
            assert resp_iou == v["iou"]

            # 5. delete address ownership
            delete_address_ownership(v["req_vasp_domain"], req_hash_addr)
            finish_delete_time = datetime.datetime.now()
            diff2 = finish_delete_time - finish_get_time
            if diff2.total_seconds() < total_time_difference:
                time.sleep(total_time_difference - diff2.total_seconds())

        except Exception as e:
            print("====== test failure ======")
            print("===== get ownership response", response)
            print("===== exception", e)
            print("====== test failure end ======")
            raise e
        count += 1


def start_test(thread_index, total_number_of_iterations):
    # generate requests
    for i in range(2):
        request_responses = request_workflow(
            thread_index, total_number_of_iterations)

        # validate responses
        response_verification(thread_index, request_responses)


def test_single_run():
    start_test(1, 1)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument('--number_of_threads')
    parser.add_argument('--total_number_of_iterations')
    args = parser.parse_args()

    threads = list()
    for index in range(int(args.number_of_threads)):
        x = threading.Thread(target=start_test, args=(
            index, int(args.total_number_of_iterations)))
        threads.append(x)
        x.start()

    for index, thread in enumerate(threads):
        thread.join()
