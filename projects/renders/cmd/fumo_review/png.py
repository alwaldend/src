"""Small, dependency-free decoder for the PNG formats used by reviews."""

import dataclasses
import os
import pathlib
import stat
import struct
import typing
import zlib

PNG_SIGNATURE = b"\x89PNG\r\n\x1a\n"
MAX_PNG_FILE_BYTES = 64 * 1024 * 1024
MAX_PNG_DIMENSION = 8192
MAX_PNG_PIXELS = 16 * 1024 * 1024
MAX_PNG_DECOMPRESSED_BYTES = 64 * 1024 * 1024
MAX_PNG_CHUNKS = 4096
_DECOMPRESS_CHUNK_BYTES = 64 * 1024


class PNGError(ValueError):
    """Raised when a PNG cannot be decoded by the review tool."""


@dataclasses.dataclass(frozen=True)
class Image:
    """An RGB image stored as tightly packed row-major bytes."""

    width: int
    height: int
    rgb: bytes

    def pixel(self, x: int, y: int) -> typing.Tuple[int, int, int]:
        """Return the RGB value at ``x``, ``y``."""
        offset = (y * self.width + x) * 3
        return (
            self.rgb[offset],
            self.rgb[offset + 1],
            self.rgb[offset + 2],
        )


@dataclasses.dataclass(frozen=True)
class _Header:
    width: int
    height: int
    color_type: int
    channels: int


def _paeth(left: int, above: int, upper_left: int) -> int:
    estimate = left + above - upper_left
    left_distance = abs(estimate - left)
    above_distance = abs(estimate - above)
    upper_left_distance = abs(estimate - upper_left)
    if (
        left_distance <= above_distance
        and left_distance <= upper_left_distance
    ):
        return left
    if above_distance <= upper_left_distance:
        return above
    return upper_left


def _unfilter_row(
    encoded: bytes,
    previous: bytes,
    bytes_per_pixel: int,
    filter_type: int,
) -> bytes:
    decoded = bytearray(len(encoded))
    for index, value in enumerate(encoded):
        left = (
            decoded[index - bytes_per_pixel] if index >= bytes_per_pixel else 0
        )
        above = previous[index] if previous else 0
        upper_left = (
            previous[index - bytes_per_pixel]
            if previous and index >= bytes_per_pixel
            else 0
        )
        if filter_type == 0:
            predictor = 0
        elif filter_type == 1:
            predictor = left
        elif filter_type == 2:
            predictor = above
        elif filter_type == 3:
            predictor = (left + above) // 2
        elif filter_type == 4:
            predictor = _paeth(left, above, upper_left)
        else:
            raise PNGError(f"unsupported PNG filter type {filter_type}")
        decoded[index] = (value + predictor) & 0xFF
    return bytes(decoded)


def _rgb_bytes(raw: bytes, color_type: int) -> bytes:
    if color_type == 2:
        return raw
    rgb = bytearray()
    if color_type == 0:
        for gray in raw:
            rgb.extend((gray, gray, gray))
        return bytes(rgb)
    if color_type == 4:
        for index in range(0, len(raw), 2):
            gray = raw[index]
            rgb.extend((gray, gray, gray))
        return bytes(rgb)
    if color_type == 6:
        for index in range(0, len(raw), 4):
            pixel_end = index + 3
            rgb.extend(raw[index:pixel_end])
        return bytes(rgb)
    raise PNGError(f"unsupported PNG color type {color_type}")


def _chunks(data: bytes) -> typing.List[typing.Tuple[bytes, bytes]]:
    if not data.startswith(PNG_SIGNATURE):
        raise PNGError("invalid PNG signature")
    cursor = len(PNG_SIGNATURE)
    chunks = []
    while cursor < len(data):
        if cursor + 12 > len(data):
            raise PNGError("truncated PNG chunk")
        length_end = cursor + 4
        type_end = cursor + 8
        length = struct.unpack(">I", data[cursor:length_end])[0]
        chunk_type = data[length_end:type_end]
        chunk_start = cursor + 8
        chunk_end = chunk_start + length
        crc_end = chunk_end + 4
        if crc_end > len(data):
            raise PNGError("truncated PNG chunk data")
        payload = data[chunk_start:chunk_end]
        expected_crc = struct.unpack(">I", data[chunk_end:crc_end])[0]
        actual_crc = zlib.crc32(chunk_type)
        actual_crc = zlib.crc32(payload, actual_crc) & 0xFFFFFFFF
        if actual_crc != expected_crc:
            name = chunk_type.decode("ascii", errors="replace")
            raise PNGError(f"CRC mismatch in {name} chunk")
        cursor = crc_end
        chunks.append((chunk_type, payload))
        if len(chunks) > MAX_PNG_CHUNKS:
            raise PNGError(f"PNG contains more than {MAX_PNG_CHUNKS} chunks")
        if chunk_type == b"IEND":
            break
    return chunks


def _validate_dimensions(width: int, height: int) -> None:
    if width <= 0 or height <= 0:
        raise PNGError("PNG dimensions must be positive")
    if width > MAX_PNG_DIMENSION or height > MAX_PNG_DIMENSION:
        raise PNGError(
            "PNG dimensions exceed limit: "
            f"{width}x{height} > {MAX_PNG_DIMENSION} per axis"
        )
    if width * height > MAX_PNG_PIXELS:
        raise PNGError(
            "PNG pixel count exceeds limit: "
            f"{width * height} > {MAX_PNG_PIXELS}"
        )


def _parse_header(chunks: typing.List[typing.Tuple[bytes, bytes]]) -> _Header:
    if not chunks or chunks[0][0] != b"IHDR" or len(chunks[0][1]) != 13:
        raise PNGError("invalid or missing IHDR chunk")
    if not any(kind == b"IEND" for kind, _ in chunks):
        raise PNGError("PNG is missing IEND")
    if sum(kind == b"IHDR" for kind, _ in chunks) != 1:
        raise PNGError("PNG must contain exactly one IHDR")
    (
        width,
        height,
        bit_depth,
        color_type,
        compression,
        filtering,
        interlace,
    ) = struct.unpack(">IIBBBBB", chunks[0][1])
    _validate_dimensions(width, height)
    if bit_depth != 8:
        raise PNGError(f"unsupported PNG bit depth {bit_depth}")
    channels_by_type = {0: 1, 2: 3, 4: 2, 6: 4}
    if color_type not in channels_by_type:
        raise PNGError(f"unsupported PNG color type {color_type}")
    if compression != 0 or filtering != 0:
        raise PNGError("unsupported PNG compression or filter method")
    if interlace != 0:
        raise PNGError("interlaced PNGs are not supported")
    return _Header(width, height, color_type, channels_by_type[color_type])


def _decompress_pending(
    decompressor: typing.Any,
    pending: bytes,
    encoded: bytearray,
    expected_size: int,
) -> None:
    while pending:
        remaining = expected_size + 1 - len(encoded)
        if remaining <= 0:
            raise PNGError(
                "PNG decompressed data exceeds expected scanline size"
            )
        encoded.extend(decompressor.decompress(pending, remaining))
        pending = decompressor.unconsumed_tail
        if len(encoded) > expected_size:
            raise PNGError(
                "PNG decompressed data exceeds expected scanline size"
            )


def _decompress_capped(compressed: bytes, expected_size: int) -> bytes:
    decompressor = zlib.decompressobj()
    encoded = bytearray()
    try:
        for offset in range(0, len(compressed), _DECOMPRESS_CHUNK_BYTES):
            chunk_end = offset + _DECOMPRESS_CHUNK_BYTES
            _decompress_pending(
                decompressor,
                compressed[offset:chunk_end],
                encoded,
                expected_size,
            )
        remaining = expected_size + 1 - len(encoded)
        if remaining <= 0 and not decompressor.eof:
            raise PNGError(
                "PNG decompressed data exceeds expected scanline size"
            )
        if remaining > 0:
            encoded.extend(decompressor.flush(remaining))
    except zlib.error as error:
        raise PNGError(f"invalid compressed PNG data: {error}") from error
    if len(encoded) > expected_size:
        raise PNGError("PNG decompressed data exceeds expected scanline size")
    if not decompressor.eof:
        raise PNGError("truncated compressed PNG data")
    if decompressor.unused_data:
        raise PNGError("compressed PNG data has trailing bytes")
    if len(encoded) != expected_size:
        raise PNGError(
            "unexpected PNG scanline size: "
            f"got {len(encoded)}, expected {expected_size}"
        )
    return bytes(encoded)


def _decode_scanlines(
    compressed: bytes, header: _Header
) -> typing.Tuple[bytes, int]:
    row_size = header.width * header.channels
    expected_size = header.height * (row_size + 1)
    if expected_size > MAX_PNG_DECOMPRESSED_BYTES:
        raise PNGError(
            "PNG decompressed data exceeds limit: "
            f"{expected_size} > {MAX_PNG_DECOMPRESSED_BYTES} bytes"
        )
    encoded = _decompress_capped(compressed, expected_size)

    raw = bytearray()
    previous = b""
    cursor = 0
    for _ in range(header.height):
        filter_type = encoded[cursor]
        cursor += 1
        row_end = cursor + row_size
        row = encoded[cursor:row_end]
        cursor += row_size
        decoded = _unfilter_row(row, previous, header.channels, filter_type)
        raw.extend(decoded)
        previous = decoded
    return bytes(raw), cursor


def decode_png_bytes(
    data: bytes, max_pixels: typing.Optional[int] = None
) -> Image:
    """Decode bounded PNG bytes within an optional remaining pixel budget."""
    if len(data) > MAX_PNG_FILE_BYTES:
        raise PNGError(
            f"PNG file exceeds limit: {len(data)} > {MAX_PNG_FILE_BYTES} bytes"
        )
    chunks = _chunks(data)
    header = _parse_header(chunks)
    pixel_count = header.width * header.height
    if max_pixels is not None and pixel_count > max_pixels:
        raise PNGError(
            "PNG pixel count exceeds remaining packet limit: "
            f"{pixel_count} > {max_pixels}"
        )
    compressed = b"".join(
        payload for kind, payload in chunks if kind == b"IDAT"
    )
    if not compressed:
        raise PNGError("PNG is missing IDAT data")
    raw, _ = _decode_scanlines(compressed, header)

    return Image(
        width=header.width,
        height=header.height,
        rgb=_rgb_bytes(raw, header.color_type),
    )


def png_dimensions_bytes(data: bytes) -> typing.Tuple[int, int]:
    """Return validated PNG dimensions without decompressing scanlines."""
    if len(data) > MAX_PNG_FILE_BYTES:
        raise PNGError(
            f"PNG file exceeds limit: {len(data)} > {MAX_PNG_FILE_BYTES} bytes"
        )
    header = _parse_header(_chunks(data))
    return header.width, header.height


def decode_png(path: pathlib.Path) -> Image:
    """Decode a bounded 8-bit, non-interlaced PNG file."""
    with path.open("rb") as source:
        metadata = os.fstat(source.fileno())
        if not stat.S_ISREG(metadata.st_mode):
            raise PNGError(f"PNG input is not a regular file: {path}")
        if metadata.st_size > MAX_PNG_FILE_BYTES:
            raise PNGError(
                "PNG file exceeds limit: "
                f"{metadata.st_size} > {MAX_PNG_FILE_BYTES} bytes"
            )
        data = source.read(MAX_PNG_FILE_BYTES + 1)
    return decode_png_bytes(data)
