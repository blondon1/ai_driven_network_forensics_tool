import joblib
import os

def load_model(model_path):
    try:
        model = joblib.load(model_path)
        return model
    except Exception as e:
        print(f"Error loading model: {e}")
        sys.exit(1)

def extract_features_from_logs(log_directory):
    features = []
    for log_file in os.listdir(log_directory):
        with open(os.path.join(log_directory, log_file), 'r', encoding="utf-8", errors='ignore') as file:
            for line in file:
                if "Packet Length" in line:
                    packet_size = int(line.split(":")[1].strip())
                    features.append([packet_size])
    return features

def main():
    model_path = 'src/ai/anomaly_detection_model.pkl'
    log_directory = 'data/logs'
    
    model = load_model(model_path)
    features = extract_features_from_logs(log_directory)
    
    for i, feature in enumerate(features):
        prediction = model.predict([feature])
        result = "Anomalous" if prediction[0] == -1 else "Normal"
        print(f"Packet {i + 1}: {result}")

if __name__ == "__main__":
    main()
