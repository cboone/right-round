#!/usr/bin/env python3
"""Split a .cast file's large output event into multiple frames by grouping every `height` lines.

Usage: python split_big_event_into_frames.py <input.cast> <out-dir>

- Finds the largest output event by length and splits its `data` by newline into blocks of `height` lines.
- Writes per-frame .cast files preserving ANSI sequences present on each line.
"""
import json, os, sys

if len(sys.argv) < 3:
    print("Usage: split_big_event_into_frames.py <input.cast> <out-dir>", file=sys.stderr)
    sys.exit(2)

inp = sys.argv[1]
outdir = sys.argv[2]

os.makedirs(outdir, exist_ok=True)

with open(inp, 'r', encoding='utf-8', errors='replace') as f:
    header = json.loads(f.readline())
    height = header.get('height', 24)

    # find the largest 'o' event (output) by data length
    largest = None
    events = []
    for line in f:
        line = line.rstrip('\n')
        if not line:
            continue
        try:
            ev = json.loads(line)
        except Exception:
            continue
        events.append(ev)
        if isinstance(ev, list) and len(ev) >= 3 and ev[1] == 'o':
            data = ev[2]
            L = len(data)
            if largest is None or L > largest[0]:
                largest = (L, ev)

    if largest is None:
        print('no output events found', file=sys.stderr)
        sys.exit(1)

    data = largest[1][2]

    # split into raw lines (preserves per-line ANSI)
    # use splitlines(keepends=False) so we can rejoin with \n
    raw_lines = data.split('\n')

    # remove trailing empty lines that may pad the frame
    while raw_lines and raw_lines[-1] == '':
        raw_lines.pop()

    if len(raw_lines) < height:
        print('largest event smaller than terminal height; nothing to split', file=sys.stderr)
        sys.exit(1)

    # group into frames of exactly `height` lines
    frames = []
    for i in range(0, len(raw_lines), height):
        block = raw_lines[i:i+height]
        if len(block) == height:
            frames.append('\n'.join(block) + '\n')

    if not frames:
        print('no full-height frames found', file=sys.stderr)
        sys.exit(1)

    # write per-frame .cast files
    for idx, block in enumerate(frames, start=1):
        outfn = os.path.join(outdir, f'frame-{idx:04d}.cast')
        with open(outfn, 'w', encoding='utf-8') as out:
            out.write(json.dumps(header) + '\n')
            # use a single output event with the original timestamp of the large chunk
            out.write(json.dumps([largest[1][0], 'o', block]) + '\n')

    print(f'Extracted {len(frames)} frames from largest event and wrote to: {outdir}')
