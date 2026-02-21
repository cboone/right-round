#!/usr/bin/env python3
"""Detect candidate frame boundaries in an asciinema v2 .cast file.

Usage: python detect_frames.py /path/to/test.cast

Heuristics used:
 - explicit CSI sequences: ESC[H, ESC[2J, hide/show cursor
 - chunks that write >= terminal height newlines
"""
import json
import sys
import re

if len(sys.argv) < 2:
    print("Usage: detect_frames.py <cast-file>", file=sys.stderr)
    sys.exit(2)

fn = sys.argv[1]
frames = []
try:
    with open(fn, 'r', encoding='utf-8', errors='replace') as f:
        header_line = f.readline()
        header = json.loads(header_line)
        height = header.get('height', 24)

        cur_frame_start_time = None
        cur_frame_start_line = None
        last_event_time = None
        event_index = 0

        for lineno, line in enumerate(f, start=2):
            line = line.strip()
            if not line:
                continue
            try:
                ev = json.loads(line)
            except Exception:
                # ignore malformed lines
                continue
            if not (isinstance(ev, list) and len(ev) >= 3):
                continue
            t, typ, data = ev[0], ev[1], ev[2]
            event_index += 1
            if cur_frame_start_time is None:
                cur_frame_start_time = t
                cur_frame_start_line = lineno

            last_event_time = t

            # explicit frame starters
            if typ == 'o' and (('\x1b[H' in data) or ('\x1b[2J' in data) or ('\x1b[?25l' in data) or ('\x1b[?25h' in data) or ('\u001b[H' in data) or ('\u001b[2J' in data)):
                frames.append({
                    'start_time': cur_frame_start_time,
                    'end_time': t,
                    'start_line': cur_frame_start_line,
                    'end_line': lineno,
                    'reason': 'explicit-csi'
                })
                cur_frame_start_time = None
                cur_frame_start_line = None
                continue

            # heuristic: writes many terminal rows
            # data is the decoded JSON string, so count actual newlines
            newline_count = data.count('\n') + data.count('\r')
            if newline_count >= max(1, height - 1):
                frames.append({
                    'start_time': cur_frame_start_time,
                    'end_time': t,
                    'start_line': cur_frame_start_line,
                    'end_line': lineno,
                    'reason': f'newline-count={newline_count}'
                })
                cur_frame_start_time = None
                cur_frame_start_line = None
                continue

        # final frame
        if cur_frame_start_time is not None:
            frames.append({
                'start_time': cur_frame_start_time,
                'end_time': last_event_time or cur_frame_start_time,
                'start_line': cur_frame_start_line,
 'end_line': lineno,
                'reason': 'eof'
            })

except FileNotFoundError:
    print(f'file not found: {fn}', file=sys.stderr)
    sys.exit(2)

# print summary
print(f"detected {len(frames)} candidate frames in {fn}\n")
for i, fr in enumerate(frames[:200], start=1):
    dur = fr['end_time'] - fr['start_time'] if fr['end_time'] is not None else 0.0
    print(f"frame {i:3d}: {fr['start_time']:.3f}s -> {fr['end_time']:.3f}s  dur={dur:.3f}s  lines={fr['start_line']}-{fr['end_line']}  reason={fr['reason']}")

# exit code 0
