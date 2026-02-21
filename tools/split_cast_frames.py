#!/usr/bin/env python3
"""Split an asciinema v2 .cast into per-frame .cast files using the same heuristics
as `detect_frames.py` (explicit CSI or newline-count >= height).

Usage: python split_cast_frames.py <input.cast> <out-dir>

Creates files: <out-dir>/frame-0001.cast, frame-0002.cast, ...
"""
import json
import os
import sys

if len(sys.argv) < 3:
    print("Usage: split_cast_frames.py <input.cast> <out-dir>", file=sys.stderr)
    sys.exit(2)

inp = sys.argv[1]
outdir = sys.argv[2]

os.makedirs(outdir, exist_ok=True)

with open(inp, 'r', encoding='utf-8', errors='replace') as f:
    header_line = f.readline()
    try:
        header = json.loads(header_line)
    except Exception as e:
        print('failed to parse header:', e, file=sys.stderr)
        sys.exit(2)
    height = header.get('height', 24)

    # read remaining lines into memory with their JSON-parsed tuples
    events = []
    for lineno, raw in enumerate(f, start=2):
        raw = raw.rstrip('\n')
        if not raw:
            continue
        try:
            ev = json.loads(raw)
        except Exception:
            # keep raw string fallback
            ev = None
        events.append((lineno, raw, ev))

# detect frames (same heuristics)
frames = []
cur_start_idx = None
for idx, (lineno, raw, ev) in enumerate(events):
    if ev is None or not isinstance(ev, list) or len(ev) < 3:
        continue
    t, typ, data = ev[0], ev[1], ev[2]
    if cur_start_idx is None:
        cur_start_idx = idx
        cur_start_time = t

    explicit = typ == 'o' and (('\x1b[H' in data) or ('\x1b[2J' in data) or ('\x1b[?25l' in data) or ('\x1b[?25h' in data) or ('\u001b[H' in data) or ('\u001b[2J' in data))
    newline_count = data.count('\n') + data.count('\r')
    if explicit or newline_count >= max(1, height - 1):
        frames.append((cur_start_idx, idx))
        cur_start_idx = None

# final
if cur_start_idx is not None:
    frames.append((cur_start_idx, len(events)-1))

if not frames:
    print('no frames detected', file=sys.stderr)
    sys.exit(1)

# write per-frame .cast files
for i, (sidx, eidx) in enumerate(frames, start=1):
    outfn = os.path.join(outdir, f'frame-{i:04d}.cast')
    with open(outfn, 'w', encoding='utf-8') as out:
        out.write(header_line)
        for _, raw, _ in events[sidx:eidx+1]:
            out.write(raw + '\n')
print(f'Wrote {len(frames)} frames to: {outdir}')
