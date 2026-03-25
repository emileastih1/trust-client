from locust import HttpUser, task, between, stats
from secrets import token_bytes
from coincurve import PublicKey
import sha3
import json
import hashlib
import sqlite3
from datetime import datetime
import os
import locust.stats

locust.stats.PERCENTILES_TO_REPORT = [0.5, 0.6, 0.7, 0.8,0.9]

PATH_TO_MY_CERT = "/test/certs_dev/client_int_1.crt"
PATH_TO_MY_KEY_FILE = "test/certs_dev/client_int_1.key"
PATH_TO_CA_BUNDLE_OR_FALSE = "test/certs_dev/bundle.crt"

now = datetime.now()

def generate_address():
    private_key = sha3.keccak_256(token_bytes(32)).digest()
    public_key = PublicKey.from_valid_secret(private_key).format(compressed=False)[1:]
    addr = sha3.keccak_256(public_key).digest()[-20:]

    return addr, private_key.hex()

def create_connection(db_file):
    conn = None
    try:
        conn = sqlite3.connect(db_file)
        cur = conn.cursor()
        cur.execute('''CREATE TABLE IF NOT EXISTS addresses (addr, addr_hash, vasp_uuid)''')
        conn.commit()
    except Exception as e:
        print(e)

    return conn

class QuickstartUser(HttpUser):
    # wait_time = between(0.5, 1)
    conn = create_connection("./loadtest.db")

    @task
    def health(self):
        response = self.client.get("/health",
                                   cert=(PATH_TO_MY_CERT, PATH_TO_MY_KEY_FILE),
                                   verify=PATH_TO_CA_BUNDLE_OR_FALSE)

    def create(self, req):
        try:
            sql = ''' INSERT INTO addresses(addr,addr_hash,vasp_uuid)
                      VALUES(?,?,?) '''
            cur = self.conn.cursor()
            cur.execute(sql, req)

            self.conn.commit()
            return cur.lastrowid
        except Exception as e:
            print(e)

    def log(self, response, id):
        log_path = './test/loadtest/results/{0}'.format(now.strftime("%m_%d_%Y-%H_%M_%S"))
        if not os.path.isdir(log_path):
            os.mkdir(log_path)
        if response.status_code >= 300:
            file = log_path + '/{1}_{2}.html'.format(now.strftime("%m_%d_%Y-%H_%M_%S"),response.status_code,id)
        else:
            file = log_path + '/{1}_{2}.html'.format(now.strftime("%m_%d_%Y-%H_%M_%S"),response.status_code,id)
        with open(file, 'wb+') as f:
            f.write(response.content)


    def get_vasp(self):

        # response = self.client.get("/vasps/0c21bcc9-b16f-442b-ab6a-503dc52c2f92",
        #                             headers={
        #                                 "x-api-token": "this-is-a-test-key-and-does-not-work-in-production"
        #                             }, cert=(PATH_TO_MY_CERT, PATH_TO_MY_KEY_FILE), verify=PATH_TO_CA_BUNDLE_OR_FALSE)
        response = self.client.get("/vasps/0c21bcc9-b16f-442b-ab6a-503dc52c2f92",
                                   headers={
                                       "x-api-token": "this-is-a-test-key-and-does-not-work-in-production"
                                   })

        self.log(response, "0c21bcc9-b16f-442b-ab6a-503dc52c2f92")


    def claim_address(self):
        # generate data
        addr, private_key = generate_address()
        addr_hash = hashlib.sha512(addr).hexdigest()
        self.create((addr.hex(), addr_hash, "991e27b2-7ebb-487f-b0dd-f3f30f03ca43"))
        response = self.client.post("/addresses",
                                    headers={
                                        "x-api-token": "this-is-a-test-key-and-does-not-work-in-production"
                                    },
                                    json={
                                        "symbol": "ETH",
                                        "address": addr_hash
                                    }, cert=(PATH_TO_MY_CERT, PATH_TO_MY_KEY_FILE), verify=PATH_TO_CA_BUNDLE_OR_FALSE)
        # response = self.client.post("/addresses",
        #                             headers={
        #                                 "x-api-token": "this-is-a-test-key-and-does-not-work-in-production"
        #                             },
        #                             json={
        #                                 "symbol": "ETH",
        #                                 "address": addr_hash
        #                             })
        self.log(response,  addr.hex())


'''
curl -vs --cacert bundle.crt --key client_int_1.key --cert client_int_1.crt https://development-rails-nlb-dev-travel-bddcaafbfdada4db.elb.us-east-1.amazonaws.com:443/_health
curl -vs --cacert bundle.crt --key client_int_1.key --cert client_int_1.crt --request POST \
  --url https://development-rails-nlb-dev-travel-bddcaafbfdada4db.elb.us-east-1.amazonaws.com:443/addresses \
  --header 'Content-Type: application/json' \
  --header 'x-api-token: this-is-a-test-key-and-does-not-work-in-production' \
  --data '{"symbol": "ETH","address": "a36e5fc7f8d526ee57cf8543e2674fd428739fa39aadb5863ed5ae640662cc089ecf3d343681c88b3ea85612de76b2800632ba1667b67ffad629419b9864e2e2"}'
locust -f  test/loadtest/locustfile.py -u 50 -r 50 -H http://localhost:3000
locust -f  test/loadtest/locustfile.py -u 20 -r 20 -H https://development-rails-nlb-dev-travel-bddcaafbfdada4db.elb.us-east-1.amazonaws.com:443 --master
locust -f  test/loadtest/locustfile.py -u 20 -r 20 -H http://development-rails-nlb-dev-travel-bddcaafbfdada4db.elb.us-east-1.amazonaws.com:3000 --master


locust -f  test/loadtest/locustfile.py -u 10 -r 10 https://development-rails-nlb-dev-travel-bddcaafbfdada4db.elb.us-east-1.amazonaws.com:443

'''