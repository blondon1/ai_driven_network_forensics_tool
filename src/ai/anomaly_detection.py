import joblib
import numpy as np
import os

model_path = 'src/ai/anomaly_detection_model.pkl'

def load_model(model_path):
    try:
        model = joblib.load(model_path)
        return model
    except Exception as e:
        print(f"Error loading model: {e}")
        return None

def is_anomalous(data_point):
    model = load_model(model_path)
    if model is None:
        return False

    data = np.array([data_point])
    prediction = model.predict(data)
    return prediction[0] == -1

# Function to check if packet size is anomalous
def check_packet_anomaly(packet_size):
    return is_anomalous([packet_size])

# To test this function independently:
if __name__ == "__main__":
    test_packet_size = 1600
    if check_packet_anomaly(test_packet_size):
        print("Anomaly detected for packet size:", test_packet_size)
    else:
        print("No anomaly detected for packet size:", test_packet_size)
