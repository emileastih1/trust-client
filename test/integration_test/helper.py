from random import randrange
from secrets import token_bytes
from coincurve import PublicKey
from sha3 import keccak_256
import hashlib


vasps = {
    "testvasp1.com": "06507e32-c5cd-4a9f-b4ab-1b2787192c8e",
    "testvasp2.com": "46d24377-ba1f-4e95-8133-946e3cfcc683",
    "testvasp3.com": "00000000-0000-0000-0000-000000000000",
    "testvasp4.com": "d0c51944-7202-4ca3-9e93-32e32d65939c",
}


address_ownership_proof_key = "addressOwnershipProof"

chain_eth = "ETHEREUM"
valid_proof_type = "ETHEREUM_CONTRACT"
valid_aux_proof_data = [
    {
        "type": "ETH_CREATE",
        "data": {
            "nonce": "123123"
        }
    },
    {
        "type": "ETH_CREATE",
        "data": {
            "nonce": "123123"
        }
    },
    {
        "type": "ETH_CREATE2",
        "data": {
            "salt": "0x14e26bf8489b5ecf8bed66029345a1fbfce9aba797dee64db4bd6403515916db",
            "init_code_hash": "0x5c5cd361ee6815781aa97a4b22dcbb49a1aed5548b57e848e34b5f3e62e3154e"
        }
    }
]

valid_aux_proof_data_2 = [
    {
        "type": "ETH_CREATE2",
        "data": {
            "salt": "0x14e26bf8489b5ecf8bed66029345a1fbfce9aba797dee64db4bd6403515916db",
            "init_code_hash": "0x5c5cd361ee6815781aa97a4b22dcbb49a1aed5548b57e848e34b5f3e62e3154e"
        }
    }
]


def get_vasp_domain_by_vasp_uuid(vasp_uuid):
    for k, v in vasps.items():
        if v == vasp_uuid:
            return k
    return None


def get_random_valid_vasp_domain():
    vasp_domains = [
        "testvasp1.com",
        "testvasp2.com",
        "testvasp3.com",
        "testvasp4.com"
    ]
    return vasp_domains[randrange(4)]


def get_vasp_uuid_by_domain_name(vasp_domain):
    return vasps[vasp_domain]


def generate_address():
    private_key = keccak_256(token_bytes(32)).digest()
    public_key = PublicKey.from_valid_secret(
        private_key).format(compressed=False)[1:]
    addr = keccak_256(public_key).digest()[-20:]
    addr_hash = hashlib.sha512(addr).hexdigest()
    return private_key, addr, addr_hash
