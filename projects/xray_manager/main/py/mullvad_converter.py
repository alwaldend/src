#!/usr/bin/env python3

import dataclasses
import urllib.parse
import enum
import requests
import argparse
import collections
import typing_extensions
import pathlib
import os
import typing
import json
import functools

T = typing.TypeVar("T")

@functools.lru_cache(1)
def iso_3166() -> typing.Mapping[str, str]:
    """
    https://gist.github.com/ssskip/5a94bfcd2835bf1dea52
    """
    return json.loads(pathlib.Path(__file__, "country-codes.json"))


@dataclasses.dataclass
class Config:
    api_url: str = "https://api.mullvad.net/app"
    token: typing.Optional[str] = None
    caddy_hostname: typing.Optional[str] = None
    clients: typing.Sequence[str] = ()
    wireguard_private_key: typing.Optional[str] = None
    wireguard_endpoint_port: int = 443
    wireguard_address: typing.Sequence[str] = ()
    aggregate_by_fields: typing.Sequence[str] = ("owned",)
    output_dir: pathlib.Path = dataclasses.field(
        default_factory=lambda: pathlib.Path("out/mullvad-converter")
    )
    subscriptions_dir: pathlib.Path = dataclasses.field(
        default_factory=lambda: pathlib.Path("xray-subscriptions")
    )
    relays_path: pathlib.Path = dataclasses.field(
        default_factory=lambda: pathlib.Path("data", "mullvad-relays.json")
    )
    wireguard_relays_aggregated: pathlib.Path = dataclasses.field(
        default_factory=lambda: pathlib.Path(
            "data", "mullvad-wireguard-relays-aggregated.json"
        )
    )
    outbounds_path: pathlib.Path = dataclasses.field(
        default_factory=lambda: pathlib.Path(
            "xray-config", "99-mullvad-outbounds.json"
        )
    )
    inbounds_path: pathlib.Path = dataclasses.field(
        default_factory=lambda: pathlib.Path(
            "xray-config", "99-mullvad-inbounds.json"
        )
    )
    rules_path: pathlib.Path = dataclasses.field(
        default_factory=lambda: pathlib.Path(
            "xray-config", "99-mullvad-rules.json"
        )
    )
    balancers_path: pathlib.Path = dataclasses.field(
        default_factory=lambda: pathlib.Path(
            "xray-config", "99-mullvad-balancers.json"
        )
    )
    caddy_path: pathlib.Path = dataclasses.field(
        default_factory=lambda: pathlib.Path(
            "caddy", "99-xray-balancers.Caddyfile"
        )
    )
    xray_socket_dir: pathlib.Path = dataclasses.field(
        default_factory=lambda: pathlib.Path("/dev/shm/xray")
    )
    http_ports_start: int = 23275

    def path(self, child: typing.Union[str, pathlib.Path]) -> pathlib.Path:
        return self.output_dir / child


def v2ray_subscription(
    protocol: str,
    id: str,
    host: str,
    query: typing.Mapping[str, typing.Any],
    comment: str,
) -> str:
    return f"{protocol}://{id}@{host}?{urllib.parse.urlencode(query)}#{comment}"


class MullvadWireguardRelay(typing.TypedDict):
    active: bool
    hostname: str
    include_in_country: bool
    ipv4_addr_in: str
    ipv6_addr_in: str
    location: str
    owned: bool
    provider: str
    public_key: str
    stboot: bool
    weight: int


class MullvadWireguardData(typing.TypedDict):
    ipv4_gateway: str
    ipv6_gateway: str
    relays: typing.Sequence[MullvadWireguardRelay]


class MullvadLocation(typing.TypedDict):
    country: str
    city: str
    lattitudude: str
    longitude: str


class MullvadWireguardRelayData(typing.TypedDict):
    locations: typing.Mapping[str, MullvadLocation]
    openvpn: typing.Mapping[str, typing.Any]
    wireguard: MullvadWireguardData


class MullvadWireguardRelaysAggregation(typing.TypedDict):
    aggregation: typing.Mapping[str, str]
    relays: typing.Sequence[MullvadWireguardRelay]


class MullvadWireguardRelaysAggregated(typing.TypedDict):
    aggregated: MullvadWireguardRelaysAggregation


class Service(str, enum.Enum):
    none = "none"
    mullvad = "mullvad"


class Protocol(str, enum.Enum):
    wireguard = "wireguard"
    http = "http"
    vless = "vless"
    freedom = "freedom"


class Transport(str, enum.Enum):
    none = "none"
    grpc = "grpc"


class OutboundTagData(typing.TypedDict):
    service: Service
    protocol: Protocol
    transport: Transport
    id: str


class InboundTagData(typing.TypedDict):
    protocol: Protocol
    transport: Transport
    id: str


@dataclasses.dataclass
class Tag(typing.Generic[T]):
    data: T

    def to_str(self) -> str:
        return json.dumps(self.data, separators=(",", ":"), sort_keys=True)

    @classmethod
    def from_str(cls, string: str) -> typing_extensions.Self:
        return cls(**json.loads(string))


BalancerTag = Tag[typing.Mapping[str, str]]
OutboundTag = Tag[OutboundTagData]
InboundTag = Tag[InboundTagData]


def relay_location(relay: MullvadWireguardRelay) -> str:
    return relay["location"].split("-", 1)[0]


@dataclasses.dataclass
class Api:
    config: Config = dataclasses.field(default_factory=Config)
    session: requests.Session = dataclasses.field(
        default_factory=requests.Session
    )

    def load_relays(
        self, *, path: typing.Optional[pathlib.Path] = None
    ) -> typing.Optional[MullvadWireguardRelayData]:
        if not path.is_file():
            return None
        with path.open("r") as file:
            return json.load(file)

    def fetch_relays(
        self,
        *,
        output_path: typing.Optional[pathlib.Path] = None,
        force: bool = False,
    ) -> MullvadWireguardRelayData:
        output_path = output_path or self.config.path(self.config.relays_path)
        cache = self.load_relays(path=output_path)
        if cache is not None and not force:
            print(f"using cache: {output_path}")
            return cache
        return self._write(output_path, self._fetch_relays())

    def wireguard_relay_outbounds(
        self,
        relay: MullvadWireguardRelay,
        secret_key: str,
        address: typing.Sequence[str],
        endpoint_port: typing.Optional[int] = None,
    ) -> list[dict[str, typing.Any]]:
        endpoint_port = endpoint_port or self.config.wireguard_endpoint_port
        if not relay["active"]:
            return []
        return [
            {
                "tag": OutboundTag(
                    data=OutboundTagData(
                        service=Service.mullvad,
                        protocol=Protocol.wireguard,
                        transport=Transport.none,
                        id=relay["hostname"],
                    )
                ).to_str(),
                "protocol": "wireguard",
                "settings": {
                    "secretKey": secret_key,
                    "address": address,
                    "peers": [
                        {
                            "endpoint": f"{relay['ipv4_addr_in']}:{endpoint_port}",
                            "publicKey": relay["public_key"],
                        }
                    ],
                },
            },
        ]

    def wireguard_relay_inbounds(
        self,
        relay: MullvadWireguardRelay,
        http_port: int,
        xray_socket_dir: pathlib.Path,
        clients: typing.Sequence[str],
    ) -> list[dict[str, typing.Any]]:
        if not relay["active"]:
            return []
        clients_list = [{"id": client} for client in clients]
        hostname = relay["hostname"]
        vless_wg_tag = InboundTag(
            data=InboundTagData(
                protocol=Protocol.vless,
                transport=Transport.grpc,
                id=OutboundTag(
                    data=OutboundTagData(
                        service=Service.mullvad,
                        protocol=Protocol.wireguard,
                        transport=Transport.none,
                        id=hostname,
                    )
                ).to_str(),
            ),
        ).to_str()
        http_tag = InboundTag(
            data=InboundTagData(
                protocol=Protocol.http,
                transport=Transport.none,
                id=OutboundTag(
                    data=OutboundTagData(
                        service=Service.mullvad,
                        protocol=Protocol.wireguard,
                        transport=Transport.none,
                        id=hostname,
                    )
                ).to_str(),
            ),
        ).to_str()
        return [
            {
                "tag": http_tag,
                "protocol": "http",
                "sniffing": {"enabled": False},
                "listen": "127.0.0.1",
                "port": http_port,
            },
            {
                "tag": vless_wg_tag,
                "protocol": "vless",
                "listen": f"{xray_socket_dir}/{vless_wg_tag}.socket,0666",
                "sniffing": {"enabled": False},
                "settings": {
                    "decryption": "none",
                    "clients": clients_list,
                },
                "streamSettings": {
                    "network": "grpc",
                    "grpcSettings": {"serviceName": vless_wg_tag},
                },
            },
        ]

    def wireguard_relay_rules(
        self,
        relay: MullvadWireguardRelay,
    ) -> list[dict[str, typing.Any]]:
        if not relay["active"]:
            return []
        hostname = relay["hostname"]
        return [
            {
                "inboundTag": [
                    InboundTag(
                        data=InboundTagData(
                            protocol=Protocol.vless,
                            transport=Transport.grpc,
                            id=OutboundTag(
                                data=OutboundTagData(
                                    service=Service.mullvad,
                                    protocol=Protocol.wireguard,
                                    transport=Transport.none,
                                    id=hostname,
                                )
                            ).to_str(),
                        )
                    ).to_str(),
                    InboundTag(
                        data=InboundTagData(
                            protocol=Protocol.http,
                            transport=Transport.none,
                            id=OutboundTag(
                                data=OutboundTagData(
                                    service=Service.mullvad,
                                    protocol=Protocol.wireguard,
                                    transport=Transport.none,
                                    id=hostname,
                                ),
                            ).to_str(),
                        ),
                    ).to_str(),
                ],
                "outboundTag": OutboundTag(
                    data=OutboundTagData(
                        service=Service.mullvad,
                        protocol=Protocol.wireguard,
                        transport=Transport.none,
                        id=hostname,
                    ),
                ).to_str(),
            }
        ]

    def wireguard_aggregation_balancers(
        self,
        aggregation: MullvadWireguardRelaysAggregation,
        balancer_tag: str,
    ) -> list[dict[str, typing.Any]]:
        balancer = {
            "tag": balancer_tag,
            "selector": [
                OutboundTag(
                    data=OutboundTagData(
                        service=Service.mullvad,
                        protocol=Protocol.wireguard,
                        transport=Transport.none,
                        id=relay["hostname"],
                    ),
                ).to_str()
                for relay in aggregation["relays"]
            ],
        }
        return [balancer]

    def wireguard_aggregation_inbounds(
        self,
        aggregation: MullvadWireguardRelaysAggregation,
        inbound_tag: str,
        xray_socket_dir: pathlib.Path,
        clients: typing.Sequence[str],
    ) -> list[dict[str, typing.Any]]:
        clients_list = [{"id": client} for client in clients]
        inbound = {
            "tag": inbound_tag,
            "protocol": "vless",
            "listen": f"{xray_socket_dir}/{inbound_tag}.socket,0666",
            "sniffing": {"enabled": False},
            "settings": {
                "decryption": "none",
                "clients": clients_list,
            },
            "streamSettings": {
                "network": "grpc",
                "grpcSettings": {"serviceName": inbound_tag},
            },
        }
        return [inbound]

    def wireguard_aggregation_rules(
        self,
        aggregation: MullvadWireguardRelaysAggregation,
        inbound_tag: str,
        balancer_tag: str,
    ) -> list[dict[str, typing.Any]]:
        rule = {"inboundTag": [inbound_tag], "balancerTag": balancer_tag}
        return [rule]

    def wireguard_balancers(
        self,
        aggregations: typing.Sequence[MullvadWireguardRelaysAggregation],
        xray_socket_dir: pathlib.Path,
        clients: typing.Sequence[str],
    ) -> dict[str, typing.Any]:
        balancers = []
        rules = []
        inbounds = []
        for aggregation in aggregations:
            balancer_tag = BalancerTag(data=aggregation["aggregation"]).to_str()
            inbound_tag = InboundTag(
                data=InboundTagData(
                    protocol=Protocol.vless,
                    transport=Transport.grpc,
                    id=balancer_tag,
                ),
            ).to_str()
            balancers.extend(
                self.wireguard_aggregation_balancers(aggregation, balancer_tag)
            )
            inbounds.extend(
                self.wireguard_aggregation_inbounds(
                    aggregation, inbound_tag, xray_socket_dir, clients
                )
            )
            rules.extend(
                self.wireguard_aggregation_rules(
                    aggregation, inbound_tag, balancer_tag
                )
            )
        return {
            "routing": {"balancers": balancers, "rules": rules},
            "inbounds": inbounds,
        }

    def write_wireguard_balancers(
        self,
        aggregated: MullvadWireguardRelaysAggregated,
        *,
        output_path: typing.Optional[pathlib.Path] = None,
        clients: typing.Sequence[str] = (),
        xray_socket_dir: typing.Optional[pathlib.Path] = None,
    ) -> None:
        output_path = output_path or self.config.path(
            self.config.balancers_path
        )
        xray_socket_dir = xray_socket_dir or self.config.xray_socket_dir
        clients = clients or self.config.clients
        if not clients:
            raise Exception("missing clients")
        result = self.wireguard_balancers(
            aggregated["aggregations"], xray_socket_dir, clients
        )
        self._write(output_path, result)

    def caddy_config(
        self,
        aggregations: typing.Sequence[MullvadWireguardRelaysAggregation],
        xray_socket_dir: pathlib.Path,
        hostname: str,
    ) -> str:
        lines = []
        for aggregation in aggregations:
            balancer_tag = BalancerTag(data=aggregation["aggregation"]).to_str()
            inbound_tag = InboundTag(
                data=InboundTagData(
                    protocol=Protocol.vless,
                    transport=Transport.grpc,
                    id=balancer_tag,
                ),
            ).to_str()
            proxy = f"""
  {hostname} {{
  @{inbound_tag} {{
    protocol grpc
    path /{inbound_tag}/*
  }}
  reverse_proxy @{inbound_tag} unix//{xray_socket_dir}/{inbound_tag}.socket {{
    flush_interval -1
    transport http {{
      versions h2c
    }}
  }}
}}
"""
            lines.append(proxy)
        return "\n".join(lines)

    def write_caddy_config(
        self,
        aggregated: MullvadWireguardRelaysAggregated,
        *,
        output_path: typing.Optional[pathlib.Path] = None,
        xray_socket_dir: typing.Optional[pathlib.Path] = None,
        hostname: typing.Optional[str] = None,
    ) -> None:
        output_path = output_path or self.config.path(self.config.caddy_path)
        hostname = hostname or self.config.caddy_hostname
        if hostname is None:
            raise Exception("missing hostname")
        xray_socket_dir = xray_socket_dir or self.config.xray_socket_dir
        result = self.caddy_config(
            aggregated["aggregations"], xray_socket_dir, hostname
        )
        self._write(output_path, result, is_json=False)

    def wireguard_balancer_subscriptions_aggregated(
        self,
        aggregated: MullvadWireguardRelaysAggregated,
        client_id: str,
        host: str,
        hostname: str,
    ) -> list[str]:
        result = []
        for aggregation in aggregated["aggregations"]:
            subscription = v2ray_subscription(
                protocol="vless",
                client_id=client_id,
                host=f"{host}:443",
                query={
                    "type": "grpc",
                    "security": "tls",
                    "sni": hostname,
                    "serviceName": BalancerTag(
                        data=aggregation["aggregation"]
                    ).to_str(),
                },
                comment=" ".join(
                    f"{key}={value}"
                    for key, value in aggregation["aggregation"].items()
                ),
            )
            result.append(subscription)
        return result

    def write_wireguard_balancer_subscriptions_aggregated(
        self,
        aggregated: MullvadWireguardRelaysAggregated,
        *,
        clients: typing.Sequence[str] = (),
        host: str = "",
        sni: str = "",
        output_path: typing.Optional[pathlib.Path] = None,
    ) -> MullvadWireguardRelaysAggregated:
        clients = clients or self.config.clients
        output_path = output_path or self.config.path(
            self.config.wireguard_relays_aggregated
        )
        return self._write(
            output_path, self.wireguard_relays_aggregated(relays=relays)
        )

    def wireguard_relays_aggregated(
        self,
        relays: MullvadWireguardRelayData,
    ) -> MullvadWireguardRelaysAggregated:
        aggr: dict[str, dict[bool, list[MullvadWireguardRelay]]] = (
            collections.defaultdict(lambda: collections.defaultdict(list))
        )
        for relay in relays["wireguard"]["relays"]:
            if not relay["active"]:
                continue
            aggr[relay_location(relay)][relay["owned"]].append(relay)
        result = []
        for location, location_relays in aggr.items():
            for owned, relays in location_relays.items():
                result.append(
                    MullvadWireguardRelaysAggregation(
                        aggregation={"loc": location, "owned": str(owned)},
                        relays=relays,
                    )
                )
        return MullvadWireguardRelaysAggregated(aggregations=result)

    def write_wireguard_relays_aggregated(
        self,
        relays: MullvadWireguardRelayData,
        *,
        output_path: typing.Optional[pathlib.Path] = None,
    ) -> MullvadWireguardRelaysAggregated:
        output_path = output_path or self.config.path(
            self.config.wireguard_relays_aggregated
        )
        return self._write(
            output_path, self.wireguard_relays_aggregated(relays=relays)
        )

    def write_wireguard_inbounds(
        self,
        relays: MullvadWireguardRelayData,
        *,
        clients: typing.Sequence[str] = (),
        http_ports_start: typing.Optional[int] = None,
        xray_socket_dir: typing.Optional[pathlib.Path] = None,
        output_path: typing.Optional[pathlib.Path] = None,
    ) -> None:
        result = []
        output_path = output_path or self.config.path(self.config.inbounds_path)
        xray_socket_dir = xray_socket_dir or self.config.xray_socket_dir
        http_ports_start = http_ports_start or self.config.http_ports_start
        clients = clients or self.config.clients
        if not clients:
            raise Exception("missing clients")

        for i, relay in enumerate(relays["wireguard"]["relays"]):
            inbounds = self.wireguard_relay_inbounds(
                relay=relay,
                xray_socket_dir=xray_socket_dir,
                http_port=http_ports_start + i + 1,
                clients=clients,
            )
            result.extend(inbounds)

        self._write(output_path, {"inbounds": result})

    def write_wireguard_outbounds(
        self,
        relays: MullvadWireguardRelayData,
        *,
        address: typing.Sequence[str] = (),
        secret_key: typing.Optional[str] = None,
        endpoint_port: typing.Optional[int] = None,
        output_path: typing.Optional[pathlib.Path] = None,
    ) -> None:
        result = []
        secret_key = secret_key or self.config.wireguard_private_key
        address = address or self.config.wireguard_address
        output_path = output_path or self.config.path(
            self.config.outbounds_path
        )
        if secret_key is None:
            raise Exception("missing wireguard secret key")
        if address is None:
            raise Exception("missing wireguard address")

        for relay in relays["wireguard"]["relays"]:
            outbounds = self.wireguard_relay_outbounds(
                relay=relay,
                secret_key=secret_key,
                address=address,
                endpoint_port=endpoint_port,
            )
            result.extend(outbounds)
        self._write(output_path, {"outbounds": result})

    def write_wireguard_rules(
        self,
        relays: MullvadWireguardRelayData,
        *,
        output_path: typing.Optional[pathlib.Path] = None,
    ) -> None:
        result = []
        output_path = output_path or self.config.path(self.config.rules_path)
        for relay in relays["wireguard"]["relays"]:
            result.extend(self.wireguard_relay_rules(relay=relay))
        self._write(output_path, {"routing": {"rules": result}})

    def _write(
        self, path: pathlib.Path, content: T, *, is_json: bool = True
    ) -> T:
        if is_json:
            content_str = json.dumps(content, indent=4, sort_keys=True)
        else:
            content_str = content
        print(f"writing '{path}'")
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(content_str + "\n")
        return content

    def _fetch_relays(self) -> MullvadWireguardRelayData:
        url = f"{self.config.api_url}/v1/relays"
        response = self.session.get(url=url)
        response.raise_for_status()
        return response.json()

    def _headers(self) -> dict[str, typing.Any]:
        if not self.config.token:
            raise Exception("missing token")
        return {"Authorization": f"Bearer {self.config.token}"}


def main() -> None:
    api = Api(
        config=Config(
            token=os.getenv("MULLVAD_XRAY_API_TOKEN"),
            wireguard_private_key=os.getenv(
                "MULLVAD_XRAY_WIREGUARD_PRIVATE_KEY"
            ),
            caddy_hostname=os.getenv("MULLVAD_XRAY_CADDY_HOSTNAME"),
            wireguard_address=tuple(
                map(
                    lambda s: s.strip(),
                    os.getenv("MULLVAD_XRAY_WIREGUARD_ADDRESS", "").split(","),
                )
            ),
            clients=tuple(
                map(
                    lambda s: s.strip(),
                    os.getenv("MULLvAD_XRAY_CLIENTS", "").split(","),
                )
            ),
        )
    )
    relays = api.fetch_relays()
    aggregated = api.write_wireguard_relays_aggregated(relays)
    api.write_wireguard_outbounds(relays)
    api.write_wireguard_inbounds(relays)
    api.write_wireguard_rules(relays)
    api.write_wireguard_balancers(aggregated)
    api.write_wireguard_balancer_subscriptions_aggregated(aggregated)
    api.write_caddy_config(aggregated)


if __name__ == "__main__":
    main()
