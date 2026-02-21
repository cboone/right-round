#!/usr/bin/env python3
"""Accumulate consecutive output events from a .cast and split into frames when
we've produced >= `height` terminal rows (counting newlines after stripping ANSI).

Usage: python split_accumulate_height.py <input.cast> <out-dir>

This handles recordings where chafa writes the screen one row per output event.
"""
import json, os, re, sys

ansi_re = re.compile(r"\x1b\[[0-9;?]*[A-Za-z]")

if len(sys.argv) < 3:
    print("Usage: split_accumulate_height.py <input.cast> <out-dir>", file=sys.stderr)
    sys.exit(2)

inp = sys.argv[1]
outdir = sys.argv[2]

os.makedirs(outdir, exist_ok=True)

with open(inp, 'r', encoding='utf-8', errors='replace') as f:
    header = json.loads(f.readline())
    height = header.get('height', 24)

    buffer_events = []
    accumulated_text = ''
    frames = []

    for line in f:
        line = line.rstrip('\n')
        if not line:
            continue
        try:
            ev = json.loads(line)
        except Exception:
            continue
        if not (isinstance(ev, list) and len(ev) >= 3):
            continue
        t, typ, data = ev[0], ev[1], ev[2]
        if typ != 'o':
            # ignore non-output for frame content but keep timing by appending
            buffer_events.append(ev)
            accumulated_text += ''
            continue

        buffer_events.append(ev)
        # strip ANSI when counting lines
        plain = ansi_re.sub('', data)
        # count actual newline characters in the chunk
        newline_count = plain.count('\n')
        # also consider chunks that don't include newline but advance cursor (rare)
        accumulated_text += plain

        # if accumulated plain text contains at least `height` full lines, cut a frame
        # compute number of complete lines in accumulated_text
        complete_lines = accumulated_text.count('\n')
        if complete_lines >= height:
            # form a single output event that concatenates the buffered output data
            concat_data = ''.join(ev[2] for ev in buffer_events if isinstance(ev, list) and len(ev) >= 3 and ev[1] == 'o')
            start_time = buffer_events[0][0]
            end_time = buffer_events[-1][0]
            frames.append((start_time, end_time, concat_data, list(buffer_events)))
            # reset
            buffer_events = []
            accumulated_text = ''

    # if leftover buffer_events, save as final partial frame
    if buffer_events:
        concat_data = ''.join(ev[2] for ev in buffer_events if isinstance(ev, list) and len(ev) >= 3 and ev[1] == 'o')
        frames.append((buffer_events[0][0], buffer_events[-1][0], concat_data, list(buffer_events)))

# write per-frame cast files
for i, (st, et, data, evs) in enumerate(frames, start=1):
    outfn = os.path.join(outdir, f'frame-{i:04d}.cast')
    with open(outfn, 'w', encoding='utf-8') as out:
        out.write(json.dumps(header) + '\n')
        # write each original event so replay timing is preserved within the frame
        for ev in evs:
            out.write(json.dumps(ev) + '\n')

print(f'Created {len(frames)} frames in {outdir} (height={height})')
