# selfpass

This is the project home of *selfpass*, the self-hosted password manager. This project is intended
to be a single-user (or **trusted** multi-user) password manager capable of encrypting/decrypting
credentials and storing them remotely through encrypted transportation, all of which is deployable
locally or to popular cloud platforms such as GCP and AWS.

It is still currently in development. However, the service is already capable of serving a gRPC based
API using mutual TLS encryption, backed by Redis and Docker. It is also capable of being deployed in
a semi-automated fashion locally and to GCP thanks to Docker.

In addition to the service there is `spc` (**s**elf**p**ass **C**LI), which is a fully fledged *selfpass* client
capable of interacting with the whole selfpass API and creating AES-CBC encrypted credentials using
a *private key* and *master password*. All of which is done using mutual TLS and an AES-CBC
encrypted config.

#### Service Roadmap

| Goal                                                                | Progress | Comment         |
| ---                                                                 | :---:    | ---             |
| Support credentials CRUD on gRPC API.                               | 100%     |                 |
| Enable server-side mutual TLS, using cfssl.                         | 100%     |                 |
| Deployable on Docker.                                               | 100%     |                 |
| Automatically deployable to GCP using docker-machine and Terraform. | 50%      | TODO: Terraform |
| Automatically deployable to AWS using docker-machine and Terraform. | 0%       |                 |

#### SPC Roadmap

| Goal                                                                   | Progress | Comment      |
| ---                                                                    | :---:    | ---          |
| Support credentials CRUD via gRPC.                                     | 80%      | TODO: Update |
| Support mutual TLS.                                                    | 100%     |              |
| Support storage of certs, PK, and host in AES-CBC encrypted config.    | 100%     |              |
| Support AES-CBC encryption of passes and OTP secrets, using MP and PK. | 100%     |              |
| Support AES-CBC encryption of local files, using MP and PK.            | 100%     |              |


#### Unplanned Goals

- Web client.
- Sensitive financial info support.
- Miscellaneous text/file encryption and storage support.
- Vault separation.

#### 3rd-party Technologies in Use (and where):
- Golang (all)
- Go-Kit (all)
- gRPC (all)
- Cobra Commander & Viper Config (spc)
- Redis (service)
- Docker (service)
- Debian (docker images and machines)
