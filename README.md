# UE-IoT-non3GPP

UE-IoT-non3GPP is an open-source implementation to provide untrusted non3GPP access to 5GCN according to 3GPP Release 15.
All the access is done via N3IWF.


**Requirements**
---

* Go 1.4.4

**Compilation**
---

```
cd ~
git clone --recursive https://github.com/my5G/my5G-core
git clone https://github.com/my5G/UE-IoT-non3GPP my5G-core/src/ue
cd my5G-core
go build -o bin/ue -x src/ue/ue.go

```

**Usage**
---

```bash
# create the ipsec link required
ip link add ipsec0 type vti local LOCAL_UE_IP remote N3IWF_IP key 5
ip link set ipsec0 up

# execute UE in background
./ue &

# trigger initial registration procedure
./src/ue/trigger_initial_request.sh
    
Options:
  --ue_addr         address in which the UE is listening to trigger procedures
  --ue_port         port in which the UE is listening to trigger procedures
  --scheme          [https|http]
  --n3iwf_address   address of the n3iwf
  --supi_or_suci    unique subscriber identifier
  --k               key
  --opc             operator code
  --ike_bind_addr   address to bind UE IKE to
```

**How to Contribute**
---

1. Fork this repo
1. Clone your fork and create a new branch: `$ git checkout https://github.com/your_username_here/UE-IoT-non3GPP -b name_for_new_branch`.
2. Make changes and test
3. Submit Pull Request with comprehensive description of changes
