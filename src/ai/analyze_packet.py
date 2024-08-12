import sys
import argparse

def main(packet_length):
    try:
        # Add your anomaly detection logic here
        print("Anomaly detected" if packet_length > 1000 else "Normal")
    except Exception as e:
        print(f"Error during processing: {e}")
        sys.exit(1)

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Analyze packet for anomalies")
    parser.add_argument('--packet-length', type=int, required=True, help='Length of the packet')
    args = parser.parse_args()

    main(args.packet_length)
