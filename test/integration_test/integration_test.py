import uuid
import requests

from test.integration_test.helper import vasps, chain_eth, valid_proof_type, valid_aux_proof_data, \
    valid_aux_proof_data_2, generate_address, address_ownership_proof_key

PATH_TO_MY_CERT = "test/certs_dev/client_int_1.crt"
PATH_TO_MY_KEY_FILE = "test/certs_dev/client_int_1.key"
PATH_TO_CA_BUNDLE_OR_FALSE = "test/certs_dev/bundle.crt"

validate_failure_message_prefix = "validation failure: "

def stage_to_test():
    return "local"


class BaseTestCase:
    def setup_method(self, method):
        self.headers = {
            'content-type': 'application/json',
        }

    def teardown_method(self, method):
        pass

    def send_request(self, method, uri, headers, body=None):
        if stage_to_test() == "development":
            base_url = "https://development-bbapi-nlb-travel-rul-503376089e6e6323.elb.us-east-1.amazonaws.com:443"
            if method == "GET":
                return requests.get(url=base_url + uri, headers=headers,
                                    cert=(PATH_TO_MY_CERT,
                                          PATH_TO_MY_KEY_FILE),
                                    verify=PATH_TO_CA_BUNDLE_OR_FALSE)
            elif method == "PUT":
                return requests.put(url=base_url + uri, headers=headers, json=body,
                                    cert=(PATH_TO_MY_CERT,
                                          PATH_TO_MY_KEY_FILE),
                                    verify=PATH_TO_CA_BUNDLE_OR_FALSE)
            elif method == "DELETE":
                return requests.delete(url=base_url + uri, headers=headers,
                                       cert=(PATH_TO_MY_CERT,
                                             PATH_TO_MY_KEY_FILE),
                                       verify=PATH_TO_CA_BUNDLE_OR_FALSE)
        elif stage_to_test() == "local":
            base_url = "http://127.0.0.1:7000"
            if method == "GET":
                return requests.get(url=base_url + uri, headers=headers)
            elif method == "PUT":
                return requests.put(url=base_url + uri, headers=headers, json=body)
            elif method == "DELETE":
                return requests.delete(url=base_url + uri, headers=headers)
        else:
            raise Exception("invalid stage")


class TestGetVASPS(BaseTestCase):
    def test_success(self):
        # given
        uri = "/v1/vasps"
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"

        # when
        response = self.send_request("GET", uri=uri, headers=self.headers)

        # then
        assert response.status_code == 200
        response_body = response.json()
        assert len(response_body["vasps"]) > 0


class TestGetVASP(BaseTestCase):
    def test_success(self):
        # given
        uri = "/v1/vasps/{0}".format(vasps["testvasp1.com"])
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"

        # when
        response = self.send_request("GET", uri=uri, headers=self.headers)

        # then
        assert response.status_code == 200
        response_body = response.json()
        vasp = response_body["vasp"]
        assert vasps["testvasp1.com"] == vasp["id"]
        assert "name" in vasp
        assert "piiEndpoint" in vasp
        assert "domain" in vasp
        assert "publicKey" in vasp
        assert "publicKeyVersion" in vasp
        assert "lei" in vasp

    def test_invalid_vasp_uuid(self):
        # given
        uri = "/v1/vasps/invalid"
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"

        # when
        response = self.send_request("GET", uri=uri, headers=self.headers)

        # then
        assert response.status_code == 400
        response_body = response.json()
        assert 3 == response_body['code']
        assert validate_failure_message_prefix + 'invalid VASP UUID' == response_body['message']

    def test_missing_requesting_vasp_should_403(self):
        # given
        uri = "/v1/vasps/{0}".format(vasps["testvasp1.com"])

        # when
        response = self.send_request("GET", uri=uri, headers=self.headers)

        # then
        assert response.status_code == 403
        response_body = response.json()
        assert 7 == response_body['code']
        assert "missing VASP identity" == response_body['message']


class TestGetAddressOwnership(BaseTestCase):
    def test_missing_chain_should_400(self):
        # given
        _, address, req_hash_addr = generate_address()
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        request_uri = "/v1/addresses/{0}".format(req_hash_addr)

        # when
        response = self.send_request(
            "GET", uri=request_uri, headers=self.headers)

        # then
        assert response.status_code == 400
        response_body = response.json()
        assert 3 == response_body['code']
        assert validate_failure_message_prefix + 'invalid chain: ' == response_body['message']

    def test_invalid_chain_should_400(self):
        # given
        _, address, req_hash_addr = generate_address()
        request_uri = "/v1/addresses/{0}?chain={1}".format(
            req_hash_addr, "invalid")
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"

        # when
        response = self.send_request(
            "GET", uri=request_uri, headers=self.headers)

        # then
        assert response.status_code == 400
        response_body = response.json()
        assert 3 == response_body['code']
        assert validate_failure_message_prefix + 'invalid chain: invalid' == response_body['message']

    def test_unauthorized_should_403(self):
        # given
        _, address, req_hash_addr = generate_address()
        request_uri = "/v1/addresses/{0}?chain={1}".format(
            req_hash_addr, chain_eth)

        # when
        response = self.send_request(
            "GET", uri=request_uri, headers=self.headers)

        # then
        assert response.status_code == 403
        response_body = response.json()
        assert 7 == response_body['code']
        assert "missing VASP identity" == response_body['message']

    def test_address_not_exist_should_404(self):
        # given
        _, address, req_hash_addr = generate_address()
        request_uri = "/v1/addresses/{0}?chain={1}".format(
            req_hash_addr, chain_eth)
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"

        # when
        response = self.send_request(
            "GET", uri=request_uri, headers=self.headers)

        # then
        assert response.status_code == 404
        response_body = response.json()
        assert 5 == response_body['code']
        assert 'address service error: owning VASP not found' == response_body['message']

    def test_invalid_address_should_400(self):
        # given
        request_uri = "/v1/addresses/{0}?chain={1}".format(
            'select * from passwords;', chain_eth)
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"

        # when
        response = self.send_request(
            "GET", uri=request_uri, headers=self.headers)

        # then
        assert response.status_code == 400
        response_body = response.json()
        assert 3 == response_body['code']
        assert validate_failure_message_prefix + 'invalid address' == response_body['message']

    def test_success(self):
        # given
        _, address, req_hash_addr = generate_address()
        uri = "/v1/addresses/{0}".format(req_hash_addr)
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        response1 = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})
        assert response1.status_code == 200
        request_uri = "/v1/addresses/{0}?chain={1}".format(
            req_hash_addr, chain_eth)

        # when
        response = self.send_request(
            "GET", uri=request_uri, headers=self.headers)

        # then
        assert response.status_code == 200
        response_body = response.json()
        assert "testvasp1.com" == response_body['vasp']['domain']
        assert response_body['addressOwnershipProof'] is None


class TestCreateAddressOwnership(BaseTestCase):
    def test_address_too_lang_should_400(self):
        # given
        uri = "/v1/addresses/{0}".format("a" * 1025)
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"

        # when
        response = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})

        # then
        assert response.status_code == 400
        response_body = response.json()
        assert 3 == response_body['code']
        assert validate_failure_message_prefix + 'invalid address' == response_body['message']

    def test_invalid_address_should_400(self):
        # given
        uri = "/v1/addresses/{0}".format(
            'select * from passwords;', chain_eth)
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"

        # when
        response = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})

        # then
        assert response.status_code == 400
        response_body = response.json()
        assert 3 == response_body['code']
        assert validate_failure_message_prefix + 'invalid address' == response_body['message']

    def test_missing_address_should_400(self):
        # given
        uri = "/v1/addresses/"
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"

        # when
        response = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})

        # then
        assert response.status_code == 400
        response_body = response.json()
        assert 3 == response_body['code']
        assert validate_failure_message_prefix + 'invalid address' == response_body['message']

    def test_missing_chain_should_400(self):
        # given
        _, address, req_hash_addr = generate_address()
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        uri = "/v1/addresses/{0}".format(req_hash_addr)

        # when
        response = self.send_request("PUT", uri=uri, headers=self.headers)

        # then
        assert response.status_code == 400
        response_body = response.json()
        assert 3 == response_body['code']
        assert validate_failure_message_prefix + 'invalid chain: ' == response_body['message']

    def test_invalid_chain_should_400(self):
        # given
        _, address, req_hash_addr = generate_address()
        uri = "/v1/addresses/{0}".format(req_hash_addr)
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"

        # when
        response = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": "invalid"})

        # then
        assert response.status_code == 400
        response_body = response.json()
        assert 3 == response_body['code']
        assert validate_failure_message_prefix + 'invalid chain: invalid' == response_body['message']

    def test_unauthorized_should_403(self):
        # given
        _, address, req_hash_addr = generate_address()
        uri = "/v1/addresses/{0}".format(req_hash_addr)

        # when
        response = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})

        # then
        assert response.status_code == 403
        response_body = response.json()
        assert 7 == response_body['code']
        assert "missing VASP identity" == response_body['message']

    def test_already_claimed_same_requesting_vasp_should_200(self):
        # given
        _, address, req_hash_addr = generate_address()
        uri = "/v1/addresses/{0}".format(req_hash_addr)
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        response1 = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})
        assert response1.status_code == 200
        response_body1 = response1.json()
        registration_id = response_body1["registrationId"]

        # when
        response = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})

        # then
        assert response.status_code == 200
        response_body = response.json()
        assert registration_id == response_body['registrationId']

    def test_already_claimed_different_requesting_vasp_should_409(self):
        # given
        _, address, req_hash_addr = generate_address()
        uri = "/v1/addresses/{0}".format(req_hash_addr)
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        response1 = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})
        assert response1.status_code == 200
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp2.com"

        # when
        response = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})

        # then
        assert response.status_code == 409
        response_body = response.json()
        assert 6 == response_body['code']
        assert "address service error: address is already claimed by another VASP" == response_body[
            'message']


class TestCreateAddressOwnershipProof(BaseTestCase):
    def test_invalid_address_should_400(self):
        # given
        _, address, req_hash_addr = generate_address()
        request_uri = "/v1/address_ownership_proofs/{0}".format(
            'select * from password;')
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        uri = "/v1/addresses/{0}".format(req_hash_addr)
        response1 = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})
        assert response1.status_code == 200
        response_body1 = response1.json()
        registration_id = response_body1["registrationId"]

        # when
        response = self.send_request("PUT", uri=request_uri, headers=self.headers, body={
            "iou": True,
            "chain": chain_eth,
            "registration_id": registration_id
        })

        # then
        assert response.status_code == 400
        response_body = response.json()
        assert 3 == response_body['code']
        assert validate_failure_message_prefix + 'invalid address' == response_body[
            'message']

    def test_invalid_chain_should_400(self):
        # given
        _, address, req_hash_addr = generate_address()
        request_uri = "/v1/address_ownership_proofs/{0}".format(req_hash_addr)
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        uri = "/v1/addresses/{0}".format(req_hash_addr)
        response1 = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})
        assert response1.status_code == 200
        response_body1 = response1.json()
        registration_id = response_body1["registrationId"]

        # when
        response = self.send_request("PUT", uri=request_uri, headers=self.headers, body={
            "iou": True,
            "chain": "invalid",
            "registration_id": registration_id
        })

        # then
        assert response.status_code == 400
        response_body = response.json()
        assert 3 == response_body['code']
        assert validate_failure_message_prefix + 'invalid chain: invalid' == response_body['message']

    def test_missing_chain_should_400(self):
        # given
        _, address, req_hash_addr = generate_address()
        request_uri = "/v1/address_ownership_proofs/{0}".format(req_hash_addr)
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        uri = "/v1/addresses/{0}".format(req_hash_addr)
        response1 = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})
        assert response1.status_code == 200
        response_body1 = response1.json()
        registration_id = response_body1["registrationId"]

        # when
        response = self.send_request("PUT", uri=request_uri, headers=self.headers, body={
            "iou": True,
            "registration_id": registration_id
        })

        # then
        assert response.status_code == 400
        response_body = response.json()
        assert 3 == response_body['code']
        assert validate_failure_message_prefix + 'invalid chain: ' == response_body['message']

    def test_invalid_iou_should_400(self):
        # given
        _, address, req_hash_addr = generate_address()
        uri = "/v1/addresses/{0}".format(req_hash_addr)
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        response1 = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})
        assert response1.status_code == 200
        response_body1 = response1.json()
        registration_id = response_body1["registrationId"]
        request_uri = "/v1/address_ownership_proofs/{0}".format(req_hash_addr)

        # when
        put_response = self.send_request("PUT", uri=request_uri, headers=self.headers, body={
            "iou": "invalid_iou",
            "chain": chain_eth,
            "registration_id": registration_id
        })

        put_response_body = put_response.json()
        print("put_response_body", put_response_body)
        assert "invalid request" == put_response_body["message"]
        assert put_response.status_code == 400

    def test_invalid_iou_as_int_should_400(self):
        # given
        _, address, req_hash_addr = generate_address()
        uri = "/v1/addresses/{0}".format(req_hash_addr)
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        response1 = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})
        assert response1.status_code == 200
        response_body1 = response1.json()
        registration_id = response_body1["registrationId"]
        request_uri = "/v1/address_ownership_proofs/{0}".format(req_hash_addr)

        # when
        put_response = self.send_request("PUT", uri=request_uri, headers=self.headers, body={
            "iou": 123,
            "chain": chain_eth,
            "registration_id": registration_id
        })

        put_response_body = put_response.json()
        print("put_response_body", put_response_body)
        assert "invalid request" == put_response_body["message"]
        assert put_response.status_code == 400

    def test_success_iou(self):
        # given
        _, address, req_hash_addr = generate_address()
        uri = "/v1/addresses/{0}".format(req_hash_addr)
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        response1 = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})
        assert response1.status_code == 200
        response_body1 = response1.json()
        registration_id = response_body1["registrationId"]
        request_uri = "/v1/address_ownership_proofs/{0}".format(req_hash_addr)

        # when
        put_response = self.send_request("PUT", uri=request_uri, headers=self.headers, body={
            "iou": True,
            "chain": chain_eth,
            "registration_id": registration_id
        })

        # then
        put_response_body = put_response.json()
        print("put_response_body", put_response_body)
        assert registration_id == put_response_body["registrationId"]
        assert put_response.status_code == 200

        # given
        get_request_uri = "/v1/addresses/{0}?chain={1}".format(
            req_hash_addr, chain_eth)

        # when
        get_response = self.send_request(
            "GET", uri=get_request_uri, headers=self.headers)

        assert get_response.status_code == 200
        get_response_body = get_response.json()
        assert "testvasp1.com" == get_response_body['vasp']['domain']
        assert get_response_body[address_ownership_proof_key] is not None
        resp_address_ownership_proof = get_response_body[address_ownership_proof_key]
        assert req_hash_addr == resp_address_ownership_proof["address"]
        assert resp_address_ownership_proof["iou"]
        assert chain_eth == resp_address_ownership_proof["chain"]
        assert resp_address_ownership_proof["signature"] is None
        assert resp_address_ownership_proof["prefix"] is None
        assert resp_address_ownership_proof["proofType"] is ''
        assert resp_address_ownership_proof["auxProofData"] == []

        # given
        put_request_uri = "/v1/address_ownership_proofs/{0}".format(
            req_hash_addr)

        put_2_response = self.send_request("PUT", uri=put_request_uri, headers=self.headers, body={
            "iou": False,
            "chain": chain_eth,
            "registration_id": registration_id,
            "prefix": "test_prefix",
            "signature": "test_signature",
            "proof_type": valid_proof_type,
            "aux_proof_data": valid_aux_proof_data
        })

        # then
        put_2_response_body = put_2_response.json()
        print(put_2_response_body)
        assert put_2_response.status_code == 200
        assert registration_id == put_2_response_body["registrationId"]

        # given
        get_request_uri = "/v1/addresses/{0}?chain={1}".format(
            req_hash_addr, chain_eth)

        # when
        get_2_response = self.send_request(
            "GET", uri=get_request_uri, headers=self.headers)

        assert get_2_response.status_code == 200
        get_2_response_body = get_2_response.json()
        assert "testvasp1.com" == get_2_response_body['vasp']['domain']
        assert get_2_response_body[address_ownership_proof_key] is not None
        get_2_resp_address_ownership_proof = get_2_response_body[address_ownership_proof_key]
        assert req_hash_addr == get_2_resp_address_ownership_proof["address"]
        assert not get_2_resp_address_ownership_proof["iou"]
        assert chain_eth == get_2_resp_address_ownership_proof["chain"]
        assert get_2_resp_address_ownership_proof["signature"] == "test_signature"
        assert get_2_resp_address_ownership_proof["prefix"] == "test_prefix"
        assert get_2_resp_address_ownership_proof["proofType"] == valid_proof_type
        assert get_2_resp_address_ownership_proof["auxProofData"] == valid_aux_proof_data

    def test_success_with_proof(self):
        # given
        _, address, req_hash_addr = generate_address()
        uri = "/v1/addresses/{0}".format(req_hash_addr)
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        response1 = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})
        assert response1.status_code == 200
        response_body1 = response1.json()
        registration_id = response_body1["registrationId"]
        request_uri = "/v1/address_ownership_proofs/{0}".format(req_hash_addr)

        # when
        response = self.send_request("PUT", uri=request_uri, headers=self.headers, body={
            "iou": False,
            "chain": chain_eth,
            "registration_id": registration_id,
            "prefix": "test_prefix",
            "signature": "test_signature",
            "proof_type": valid_proof_type,
            "aux_proof_data": valid_aux_proof_data
        })

        # then
        assert response.status_code == 200
        response_body = response.json()
        assert registration_id == response_body["registrationId"]

        # given
        get_request_uri = "/v1/addresses/{0}?chain={1}".format(
            req_hash_addr, chain_eth)

        # when
        get_response = self.send_request(
            "GET", uri=get_request_uri, headers=self.headers)

        assert get_response.status_code == 200
        get_response_body = get_response.json()
        assert "testvasp1.com" == get_response_body['vasp']['domain']
        assert get_response_body[address_ownership_proof_key] is not None
        resp_address_ownership_proof = get_response_body[address_ownership_proof_key]
        assert req_hash_addr == resp_address_ownership_proof["address"]
        assert not resp_address_ownership_proof["iou"]
        assert chain_eth == resp_address_ownership_proof["chain"]
        assert resp_address_ownership_proof["signature"] == "test_signature"
        assert resp_address_ownership_proof["prefix"] == "test_prefix"
        assert resp_address_ownership_proof["proofType"] == valid_proof_type
        assert resp_address_ownership_proof["auxProofData"] == valid_aux_proof_data

        # when
        response2 = self.send_request("PUT", uri=request_uri, headers=self.headers, body={
            "iou": False,
            "chain": chain_eth,
            "registration_id": registration_id,
            "prefix": "test_prefix",
            "signature": "test_signature",
            "proof_type": valid_proof_type,
            "aux_proof_data": valid_aux_proof_data_2
        })

        # then
        assert response2.status_code == 200
        response_body2 = response2.json()
        assert registration_id == response_body2["registrationId"]

        # given
        get_request_uri = "/v1/addresses/{0}?chain={1}".format(
            req_hash_addr, chain_eth)

        # when
        get_response = self.send_request(
            "GET", uri=get_request_uri, headers=self.headers)

        assert get_response.status_code == 200
        get_response_body = get_response.json()
        assert "testvasp1.com" == get_response_body['vasp']['domain']
        assert get_response_body[address_ownership_proof_key] is not None
        resp_address_ownership_proof = get_response_body[address_ownership_proof_key]
        assert req_hash_addr == resp_address_ownership_proof["address"]
        assert not resp_address_ownership_proof["iou"]
        assert chain_eth == resp_address_ownership_proof["chain"]
        assert resp_address_ownership_proof["signature"] == "test_signature"
        assert resp_address_ownership_proof["prefix"] == "test_prefix"
        assert resp_address_ownership_proof["proofType"] == valid_proof_type
        assert resp_address_ownership_proof["auxProofData"] == valid_aux_proof_data_2

    def test_invalid_proof_type(self):
        # given
        _, address, req_hash_addr = generate_address()
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        uri = "/v1/addresses/{0}".format(req_hash_addr)
        response1 = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})
        assert response1.status_code == 200
        response_body1 = response1.json()
        registration_id = response_body1["registrationId"]
        data = {
            "iou": False,
            "chain": chain_eth,
            "registration_id": registration_id,
            "signature": "test_signature",
            "prefix": "test_prefix",
            "proof_type": "invalid",
            "aux_proof_data": valid_aux_proof_data
        }
        request_uri = "/v1/address_ownership_proofs/{0}".format(req_hash_addr)

        # when
        response = self.send_request(
            "PUT", uri=request_uri, headers=self.headers, body=data)

        # then
        assert response.status_code == 400
        response_body = response.json()
        assert 3 == response_body['code']
        assert validate_failure_message_prefix + 'invalid proof type' ==  response_body["message"]

    def test_invalid_aux_proof_data(self):
        # given
        _, address, req_hash_addr = generate_address()
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        uri = "/v1/addresses/{0}".format(req_hash_addr)
        response1 = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})
        assert response1.status_code == 200
        response_body1 = response1.json()
        registration_id = response_body1["registrationId"]
        data = {
            "iou": False,
            "chain": chain_eth,
            "registration_id": registration_id,
            "signature": "test_signature",
            "prefix": "test_prefix",
            "proof_type": valid_proof_type,
            "aux_proof_data": [{
                "type":"ETH_CREATE",
                "data": {
                    "<img src=a onerror=alert(1)>": "alert"
                }
            }]
        }
        request_uri = "/v1/address_ownership_proofs/{0}".format(req_hash_addr)

        # when
        response = self.send_request(
            "PUT", uri=request_uri, headers=self.headers, body=data)

        # then
        assert response.status_code == 400
        response_body = response.json()
        assert 3 == response_body['code']
        assert validate_failure_message_prefix + 'invalid aux proof data' == response_body["message"]

    def test_invalid_aux_proof_data_invalid_eth_create2_missing_salt(self):
        # given
        _, address, req_hash_addr = generate_address()
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        uri = "/v1/addresses/{0}".format(req_hash_addr)
        response1 = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})
        assert response1.status_code == 200
        response_body1 = response1.json()
        registration_id = response_body1["registrationId"]
        data = {
            "iou": False,
            "chain": chain_eth,
            "registration_id": registration_id,
            "signature": "test_signature",
            "prefix": "test_prefix",
            "proof_type": valid_proof_type,
            "aux_proof_data": [
                {
                    "type": "ETH_CREATE2",
                    "data": {
                        "init_code_hash": "0x5c5cd361ee6815781aa97a4b22dcbb49a1aed5548b57e848e34b5f3e62e3154e"
                    }
                }
            ]
        }
        request_uri = "/v1/address_ownership_proofs/{0}".format(req_hash_addr)

        # when
        response = self.send_request(
            "PUT", uri=request_uri, headers=self.headers, body=data)

        # then
        assert response.status_code == 400
        response_body = response.json()
        assert 3 == response_body['code']
        assert validate_failure_message_prefix + 'address service error: invalid proof type and data combination' == response_body["message"]


    def test_invalid_aux_proof_data_invalid_eth_create2_missing_salt(self):
        # given
        _, address, req_hash_addr = generate_address()
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        uri = "/v1/addresses/{0}".format(req_hash_addr)
        response1 = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})
        assert response1.status_code == 200
        response_body1 = response1.json()
        registration_id = response_body1["registrationId"]
        data = {
            "iou": False,
            "chain": chain_eth,
            "registration_id": registration_id,
            "signature": "test_signature",
            "prefix": "test_prefix",
            "proof_type": valid_proof_type,
            "aux_proof_data": [
                {
                    "type": "ETH_CREATE2",
                    "data": {
                        "salt": "14e26bf8489b5ecf8bed66029345a1fbfce9aba797dee64db4bd6403515916db",
                        "init_code_hash": "5c5cd361ee6815781aa97a4b22dcbb49a1aed5548b57e848e34b5f3e62e3154e"
                    }
                }
            ]
        }
        request_uri = "/v1/address_ownership_proofs/{0}".format(req_hash_addr)

        # when
        response = self.send_request(
            "PUT", uri=request_uri, headers=self.headers, body=data)

        # then
        assert response.status_code == 400
        response_body = response.json()
        assert 3 == response_body['code']
        assert validate_failure_message_prefix + 'invalid aux proof data' == response_body["message"]

    def test_invalid_proof_type_should_400(self):
        # given
        _, address, req_hash_addr = generate_address()
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        uri = "/v1/addresses/{0}".format(req_hash_addr)
        response1 = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})
        assert response1.status_code == 200
        response_body1 = response1.json()
        registration_id = response_body1["registrationId"]
        data = {
            "iou": False,
            "chain": chain_eth,
            "registration_id": registration_id,
            "signature": "test_signature",
            "prefix": "test_prefix",
            "proof_type": "<img src=a onerror=alert(1)>",
            "aux_proof_data": [{
                "type":"ETH_CREATE",
                "data": {
                    "params": "test"
                }
            }]
        }
        request_uri = "/v1/address_ownership_proofs/{0}".format(req_hash_addr)

        # when
        response = self.send_request(
            "PUT", uri=request_uri, headers=self.headers, body=data)

        # then
        assert response.status_code == 400
        response_body = response.json()
        assert 3 == response_body['code']
        assert validate_failure_message_prefix + 'invalid proof type' == response_body["message"]

    def test_invalid_signature_should_400(self):
        # given
        _, address, req_hash_addr = generate_address()
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        uri = "/v1/addresses/{0}".format(req_hash_addr)
        response1 = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})
        assert response1.status_code == 200
        response_body1 = response1.json()
        registration_id = response_body1["registrationId"]
        data = {
            "iou": False,
            "chain": chain_eth,
            "registration_id": registration_id,
            "signature": "<img src=a onerror=alert(1)>",
            "prefix": "test_prefix",
            "proof_type": valid_proof_type,
            "aux_proof_data": [{
                "type":"ETH_CREATE",
                "data": {
                    "params": "test"
                }
            }]
        }
        request_uri = "/v1/address_ownership_proofs/{0}".format(req_hash_addr)

        # when
        response = self.send_request(
            "PUT", uri=request_uri, headers=self.headers, body=data)

        # then
        assert response.status_code == 400
        response_body = response.json()
        assert 3 == response_body['code']
        assert validate_failure_message_prefix + 'invalid signature' == response_body["message"]

    def test_invalid_aux_proof_data_invalid_value(self):
        # given
        _, address, req_hash_addr = generate_address()
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        uri = "/v1/addresses/{0}".format(req_hash_addr)
        response1 = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})
        assert response1.status_code == 200
        response_body1 = response1.json()
        registration_id = response_body1["registrationId"]
        data = {
            "iou": False,
            "chain": chain_eth,
            "registration_id": registration_id,
            "signature": "test_signature",
            "prefix": "test_prefix",
            "proof_type": valid_proof_type,
            "aux_proof_data": [{
                "type":"ETH_CREATE",
                "data": {
                    "params": "<img src=a onerror=alert(1)>"
                }
            }]
        }
        request_uri = "/v1/address_ownership_proofs/{0}".format(req_hash_addr)

        # when
        response = self.send_request(
            "PUT", uri=request_uri, headers=self.headers, body=data)

        # then
        assert response.status_code == 400
        response_body = response.json()
        assert 3 == response_body['code']
        assert validate_failure_message_prefix + 'invalid aux proof data' == response_body["message"]

    def test_invalid_aux_proof_data_with_invalid_key(self):
        # given
        _, address, req_hash_addr = generate_address()
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        uri = "/v1/addresses/{0}".format(req_hash_addr)
        response1 = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})
        assert response1.status_code == 200
        response_body1 = response1.json()
        registration_id = response_body1["registrationId"]
        data = {
            "iou": False,
            "chain": chain_eth,
            "registration_id": registration_id,
            "signature": "test_signature",
            "prefix": "test_prefix",
            "proof_type": valid_proof_type,
            "aux_proof_data": [
                {
                    "type": "<img src=a onerror=alert(1)>",
                    "data": {
                        "params2": "test_params2"
                    }
                }
            ]
        }
        request_uri = "/v1/address_ownership_proofs/{0}".format(req_hash_addr)

        # when
        response = self.send_request(
            "PUT", uri=request_uri, headers=self.headers, body=data)

        # then
        assert response.status_code == 400
        response_body = response.json()
        assert 3 == response_body['code']
        assert validate_failure_message_prefix + 'invalid aux proof data' == response_body["message"]

    def test_invalid_aux_proof_data_missing_type(self):
        # given
        _, address, req_hash_addr = generate_address()
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        uri = "/v1/addresses/{0}".format(req_hash_addr)
        response1 = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})
        assert response1.status_code == 200
        response_body1 = response1.json()
        registration_id = response_body1["registrationId"]
        data = {
            "iou": False,
            "chain": chain_eth,
            "registration_id": registration_id,
            "signature": "test_signature",
            "prefix": "test_prefix",
            "proof_type": valid_proof_type,
            "aux_proof_data": [
                {
                    "data": {
                        "params2": "test_params2"
                    }
                }
            ]
        }
        request_uri = "/v1/address_ownership_proofs/{0}".format(req_hash_addr)

        # when
        response = self.send_request(
            "PUT", uri=request_uri, headers=self.headers, body=data)

        # then
        assert response.status_code == 400
        response_body = response.json()
        assert 3 == response_body['code']
        assert validate_failure_message_prefix + 'invalid aux proof data' == response_body["message"]

    def test_invalid_aux_proof_data_with_missing_data(self):
        # given
        _, address, req_hash_addr = generate_address()
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        uri = "/v1/addresses/{0}".format(req_hash_addr)
        response1 = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})
        assert response1.status_code == 200
        response_body1 = response1.json()
        registration_id = response_body1["registrationId"]
        data = {
            "iou": False,
            "chain": chain_eth,
            "registration_id": registration_id,
            "signature": "test_signature",
            "prefix": "test_prefix",
            "proof_type": valid_proof_type,
            "aux_proof_data": [
                {
                    "type": "ETH_CREATE",
                }
            ]
        }
        request_uri = "/v1/address_ownership_proofs/{0}".format(req_hash_addr)

        # when
        response = self.send_request(
            "PUT", uri=request_uri, headers=self.headers, body=data)

        # then
        assert response.status_code == 400
        response_body = response.json()
        assert 3 == response_body['code']
        assert validate_failure_message_prefix + 'invalid aux proof data' == response_body["message"]

    def test_iou_true_proof_provided_should_400(self):
        # given
        _, address, req_hash_addr = generate_address()
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        uri = "/v1/addresses/{0}".format(req_hash_addr)
        response1 = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})
        assert response1.status_code == 200
        response_body1 = response1.json()
        registration_id = response_body1["registrationId"]
        data = {
            "iou": True,
            "chain": chain_eth,
            "registration_id": registration_id,
            "signature": "test_signature",
            "prefix": "test_prefix",
            "proof_type": valid_proof_type,
            "aux_proof_data": valid_aux_proof_data
        }
        request_uri = "/v1/address_ownership_proofs/{0}".format(req_hash_addr)

        # when
        response = self.send_request(
            "PUT", uri=request_uri, headers=self.headers, body=data)

        # then
        assert response.status_code == 400
        response_body = response.json()
        assert 3 == response_body['code']
        assert validate_failure_message_prefix + "address service error: can not set proof data when IOU is true" == response_body[
            "message"]

    def test_iou_true_sig_provided_should_400(self):
        # given
        _, address, req_hash_addr = generate_address()
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        uri = "/v1/addresses/{0}".format(req_hash_addr)
        response1 = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})
        assert response1.status_code == 200
        response_body1 = response1.json()
        registration_id = response_body1["registrationId"]
        data = {
            "iou": True,
            "chain": chain_eth,
            "registration_id": registration_id,
            "signature": "test_signature",
        }
        request_uri = "/v1/address_ownership_proofs/{0}".format(req_hash_addr)

        # when
        response = self.send_request(
            "PUT", uri=request_uri, headers=self.headers, body=data)

        # then
        assert response.status_code == 400
        response_body = response.json()
        assert 3 == response_body['code']
        assert validate_failure_message_prefix + "address service error: can not set proof data when IOU is true" == response_body[
            "message"]

    def test_iou_true_prefix_provided_should_400(self):
        # given
        _, address, req_hash_addr = generate_address()
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        uri = "/v1/addresses/{0}".format(req_hash_addr)
        response1 = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})
        assert response1.status_code == 200
        response_body1 = response1.json()
        registration_id = response_body1["registrationId"]
        data = {
            "iou": True,
            "chain": chain_eth,
            "registration_id": registration_id,
            "prefix": "test_prefix",
        }
        request_uri = "/v1/address_ownership_proofs/{0}".format(req_hash_addr)

        # when
        response = self.send_request(
            "PUT", uri=request_uri, headers=self.headers, body=data)

        # then
        assert response.status_code == 400
        response_body = response.json()
        assert 3 == response_body['code']
        assert validate_failure_message_prefix + "address service error: can not set proof data when IOU is true" == response_body[
            "message"]

    def test_iou_true_prooftype_provided_should_400(self):
        # given
        _, address, req_hash_addr = generate_address()
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        uri = "/v1/addresses/{0}".format(req_hash_addr)
        response1 = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})
        assert response1.status_code == 200
        response_body1 = response1.json()
        registration_id = response_body1["registrationId"]
        data = {
            "iou": True,
            "chain": chain_eth,
            "registration_id": registration_id,
            "proof_type": valid_proof_type,
        }
        request_uri = "/v1/address_ownership_proofs/{0}".format(req_hash_addr)

        # when
        response = self.send_request(
            "PUT", uri=request_uri, headers=self.headers, body=data)

        # then
        assert response.status_code == 400
        response_body = response.json()
        assert 3 == response_body['code']
        assert validate_failure_message_prefix + "address service error: can not set proof data when IOU is true" == response_body[
            "message"]

    def test_iou_true_aux_proof_data_provided_should_400(self):
        # given
        _, address, req_hash_addr = generate_address()
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        uri = "/v1/addresses/{0}".format(req_hash_addr)
        response1 = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})
        assert response1.status_code == 200
        response_body1 = response1.json()
        registration_id = response_body1["registrationId"]
        data = {
            "iou": True,
            "chain": chain_eth,
            "registration_id": registration_id,
            "aux_proof_data": valid_aux_proof_data
        }
        request_uri = "/v1/address_ownership_proofs/{0}".format(req_hash_addr)

        # when
        response = self.send_request(
            "PUT", uri=request_uri, headers=self.headers, body=data)

        # then
        assert response.status_code == 400
        response_body = response.json()
        assert 3 == response_body['code']
        assert validate_failure_message_prefix + "address service error: can not set proof data when IOU is true" == response_body[
            "message"]

    def test_iou_false_and_missing_prefix_should_400(self):
        # given
        _, address, req_hash_addr = generate_address()
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        uri = "/v1/addresses/{0}".format(req_hash_addr)
        response1 = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})
        assert response1.status_code == 200
        response_body1 = response1.json()
        registration_id = response_body1["registrationId"]
        data = {
            "iou": False,
            "chain": chain_eth,
            "registration_id": registration_id,
            "signature": "test_signature",
            "proof_type": valid_proof_type,
            "aux_proof_data": valid_aux_proof_data
        }
        request_uri = "/v1/address_ownership_proofs/{0}".format(req_hash_addr)

        # when
        response = self.send_request(
            "PUT", uri=request_uri, headers=self.headers, body=data)

        # then
        assert response.status_code == 400
        response_body = response.json()
        assert 3 == response_body['code']
        assert validate_failure_message_prefix + "address service error: invalid proof, prefix is empty" == response_body[
            "message"]

    def test_iou_false_and_missing_signature_should_400(self):
        # given
        _, address, req_hash_addr = generate_address()
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        uri = "/v1/addresses/{0}".format(req_hash_addr)
        response1 = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})
        assert response1.status_code == 200
        response_body1 = response1.json()
        registration_id = response_body1["registrationId"]
        data = {
            "iou": False,
            "chain": chain_eth,
            "registration_id": registration_id,
            "prefix": "test_prefix",
            "proof_type": valid_proof_type,
            "aux_proof_data": valid_aux_proof_data
        }
        request_uri = "/v1/address_ownership_proofs/{0}".format(req_hash_addr)

        # when
        response = self.send_request(
            "PUT", uri=request_uri, headers=self.headers, body=data)

        # then
        assert response.status_code == 400
        response_body = response.json()
        assert 3 == response_body['code']
        assert validate_failure_message_prefix + "address service error: invalid proof, signature is empty" == response_body[
            "message"]

    def test_iou_false_and_missing_prooftype_should_400(self):
        # given
        _, address, req_hash_addr = generate_address()
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        uri = "/v1/addresses/{0}".format(req_hash_addr)
        response1 = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})
        assert response1.status_code == 200
        response_body1 = response1.json()
        registration_id = response_body1["registrationId"]
        data = {
            "iou": False,
            "chain": chain_eth,
            "registration_id": registration_id,
            "prefix": "test_prefix",
            "signature": "test_signature",
            "aux_proof_data": valid_aux_proof_data
        }
        request_uri = "/v1/address_ownership_proofs/{0}".format(req_hash_addr)

        # when
        response = self.send_request(
            "PUT", uri=request_uri, headers=self.headers, body=data)

        # then
        assert response.status_code == 400
        response_body = response.json()
        assert 3 == response_body['code']
        assert validate_failure_message_prefix + "address service error: proof type must not be null when IOU is false" == response_body[
            "message"]

    def test_iou_false_and_missing_aux_proof_data_should_200(self):
        # given
        _, address, req_hash_addr = generate_address()
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        uri = "/v1/addresses/{0}".format(req_hash_addr)
        response1 = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})
        assert response1.status_code == 200
        response_body1 = response1.json()
        registration_id = response_body1["registrationId"]
        data = {
            "iou": False,
            "chain": chain_eth,
            "registration_id": registration_id,
            "prefix": "test_prefix",
            "signature": "test_signature",
            "proof_type": "ETHEREUM_EOA"
        }
        request_uri = "/v1/address_ownership_proofs/{0}".format(req_hash_addr)

        # when
        response = self.send_request(
            "PUT", uri=request_uri, headers=self.headers, body=data)

        # then
        assert response.status_code == 200

    def test_missing_registration_id_should_400(self):
        # given
        _, address, req_hash_addr = generate_address()
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        request_uri = "/v1/address_ownership_proofs/{0}".format(req_hash_addr)
        data = {
            "iou": True,
            "chain": chain_eth,
        }

        # when
        response = self.send_request(
            "PUT", uri=request_uri, headers=self.headers, body=data)

        # then
        assert response.status_code == 400
        response_body = response.json()
        assert 3 == response_body['code']
        assert validate_failure_message_prefix + "invalid registration id" == response_body[
            "message"]

    def test_invalid_registration_id_should_400(self):
        # given
        _, address, req_hash_addr = generate_address()
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        request_uri = "/v1/address_ownership_proofs/{0}".format(req_hash_addr)
        data = {
            "iou": True,
            "chain": chain_eth,
            "registration_id": str(uuid.uuid4()),
            "proof_type": valid_proof_type,
            "aux_proof_data": valid_aux_proof_data
        }

        # when
        response = self.send_request(
            "PUT", uri=request_uri, headers=self.headers, body=data)

        # then
        assert response.status_code == 400
        response_body = response.json()
        assert 3 == response_body['code']
        assert validate_failure_message_prefix + "address service error: updating proof for an unclaimed address" == response_body[
            "message"]

    def test_provide_proof_owned_by_other_vasp_should_403(self):
        # given
        _, address, req_hash_addr = generate_address()
        uri = "/v1/addresses/{0}".format(req_hash_addr)
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        response1 = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})
        assert response1.status_code == 200
        response_body1 = response1.json()
        registration_id = response_body1["registrationId"]
        request_uri = "/v1/address_ownership_proofs/{0}".format(req_hash_addr)
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp2.com"

        # when
        response = self.send_request("PUT", uri=request_uri, headers=self.headers, body={
            "iou": True,
            "chain": chain_eth,
            "registration_id": registration_id,
            "proof_type": valid_proof_type,
            "aux_proof_data": valid_aux_proof_data
        })

        # then
        assert response.status_code == 403
        response_body = response.json()
        assert 7 == response_body['code']
        assert "address service error: requesting VASP and owner VASP does not match" == response_body[
            "message"]


class TestDeleteAddressOwnership(BaseTestCase):
    def test_success(self):
        # given
        _, address, req_hash_addr = generate_address()
        uri = "/v1/addresses/{0}".format(req_hash_addr)
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        response1 = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})
        assert response1.status_code == 200
        request_uri = "/v1/addresses/{0}?chain={1}".format(
            req_hash_addr, chain_eth)

        # when
        response = self.send_request(
            "DELETE", uri=request_uri, headers=self.headers)

        # then
        assert response.status_code == 200

    def test_invalid_address_should_400(self):
        # given
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        request_uri = "/v1/addresses/{0}?chain={1}".format("select * from users;", chain_eth)

        # when
        response = self.send_request(
            "DELETE", uri=request_uri, headers=self.headers)

        # then
        assert response.status_code == 400
        response_body = response.json()
        assert 3 == response_body['code']
        assert validate_failure_message_prefix + 'invalid address' == response_body["message"]

    def test_delete_no_chain_should_400(self):
        # given
        _, address, req_hash_addr = generate_address()
        uri = "/v1/addresses/{0}".format(req_hash_addr)
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        response1 = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})
        assert response1.status_code == 200
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        request_uri = "/v1/addresses/{0}".format(req_hash_addr)

        # when
        response = self.send_request(
            "DELETE", uri=request_uri, headers=self.headers)

        # then
        assert response.status_code == 400
        response_body = response.json()
        assert 3 == response_body['code']
        assert validate_failure_message_prefix + 'invalid chain: ' == response_body["message"]

    def test_delete_vasp_without_permission_should_403(self):
        # given
        _, address, req_hash_addr = generate_address()
        uri = "/v1/addresses/{0}".format(req_hash_addr)
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        response1 = self.send_request(
            "PUT", uri=uri, headers=self.headers, body={"chain": chain_eth})
        assert response1.status_code == 200
        request_uri = "/v1/addresses/{0}?chain={1}".format(
            req_hash_addr, chain_eth)
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp2.com"

        # when
        response = self.send_request(
            "DELETE", uri=request_uri, headers=self.headers)

        # then
        assert response.status_code == 403
        response_body = response.json()
        assert 7 == response_body['code']
        assert "address service error: can not delete address ownership that is not owned by the current requesting VASP" == response_body[
            "message"]

    def test_invalid_chain_should_400(self):
        # given
        _, address, req_hash_addr = generate_address()
        self.headers['grpc-metadata-x-vasp-dns-name'] = "testvasp1.com"
        request_uri = "/v1/addresses/{0}?chain={1}".format(
            req_hash_addr, "invalid")

        # when
        response = self.send_request(
            "DELETE", uri=request_uri, headers=self.headers)

        # then
        assert response.status_code == 400
        response_body = response.json()
        assert 3 == response_body['code']
        assert validate_failure_message_prefix + 'invalid chain: invalid' == response_body["message"]
