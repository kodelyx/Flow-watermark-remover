import argparse
import subprocess
import sys
import os

def main():
    parser = argparse.ArgumentParser(description="Add brand watermark to video")
    parser.add_argument("-i", "--input", required=True, help="Input video file")
    parser.add_argument("-o", "--output", required=True, help="Output video file")
    parser.add_argument("-w", "--watermark", default="Akash Digital Marketing", help="Watermark text")
    args = parser.parse_args()

    if not os.path.exists(args.input):
        print(f"Error: Input file {args.input} does not exist.")
        sys.exit(1)

    # FFmpeg command to overlay text watermark centered in the video
    # drawtext is centered: x=(w-text_w)/2, y=(h-text_h)/2
    cmd = [
        "ffmpeg", "-y", "-i", args.input,
        "-vf", f"drawtext=text='{args.watermark}':x=(w-text_w)/2:y=(h-text_h)/2:fontsize=36:fontcolor=white@0.4:box=1:boxcolor=black@0.3:boxborderw=10",
        "-c:a", "copy", args.output
    ]

    print(f"Adding brand watermark to {args.input}...")
    try:
        subprocess.run(cmd, check=True)
        print(f"Success! Watermarked video saved to {args.output}")
    except subprocess.CalledProcessError as e:
        print(f"Error executing FFmpeg: {e}")
        sys.exit(1)
    except FileNotFoundError:
        print("Error: FFmpeg is not installed or not in PATH.")
        sys.exit(1)

if __name__ == "__main__":
    main()
