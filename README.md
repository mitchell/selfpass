# selfpass

This is the project home of *selfpass*, the self-hosted password manager. This project is
a single-user password manager capable of encrypting/decrypting credentials and storing them
remotely through encrypted transportation. All of which is deployable locally or to popular cloud
platforms such as GCP and AWS.

It is still currently in development. However, the server is already capable of serving a gRPC
based API using mutual TLS encryption, backed by BoltDB. It is also capable
of being deployed in a semi-automated fashion locally and to GCP thanks to Docker.

**Server Roadmap**

| Goal                                                                | Progress | Comment         |
| ---                                                                 | :---:    | ---             |
| Support credentials CRUD on gRPC API.                               | 100%     |                 |
| Enable server-side mutual TLS, using cfssl.                         | 100%     |                 |
| Deployable on Docker.                                               | 100%     |                 |
| Automatically deployable to GCP using docker-machine and Terraform. | 50%      | TODO: Terraform |
| Automatically deployable to AWS using docker-machine and Terraform. | 0%       |                 |

## sp CLI

In addition to the server there is `sp`, which is a fully fledged *selfpass* client capable of
interacting with the whole selfpass API and creating AES-CBC encrypted credentials using a *private
key* and *master password*. All of which is done using mutual TLS and an AES-GCM encrypted config.

**CLI Roadmap**

| Goal                                                                   | Progress | Comment |
| ---                                                                    | :---:    | ---     |
| Support mutual TLS.                                                    | 100%     |         |
| Support credentials CRUD via gRPC.                                     | 100%     |         |
| Support storage of certs, PK, and host in AES-GCM encrypted config.    | 100%     |         |
| Support AES-CBC encryption of passes and OTP secrets, using MP and PK. | 100%     |         |
| Support AES-GCM encryption of local files, using MP and PK.            | 100%     |         |

## Client

The newest addition to the *selfpass* project is the client built using Flutter, which makes it
capable of targeting to iOS, Android, and Desktop. It supports all the same features as the CLI tool
using GUIs, with all the same safety and encryption as the CLI.

| Goal                                                                     | Progress | Comment   |
| ---                                                                      | :---:    | ---       |
| Support mutual TLS.                                                      | 100%     |           |
| Support credentials CRUD via gRPC.                                       | 25%      | TODO: CUD |
| Support storage of certs, PK, and host in shared preferences, encrypted. | 100%     |           |
| Support AES-CBC encryption of passes and OTP secrets, using MP and PK.   | 100%     |           |

## Other Info

**Unplanned Goals**

- Sensitive financial info support.
- Miscellaneous text/file encryption and storage support.
- Vault separation.

**Architectural 3rd-party Technologies in Use (and where)**

- Golang: services, sp, & protobuf
- Dart: client & protobuf
- Flutter: client
- Go-Kit: services
- gRPC/Protobuf: all
- Cobra Commander & Viper Config: sp
- BoltDB/Redis: services
- Docker: services
- Debian: docker images & machines
