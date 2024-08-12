from sklearn.ensemble import IsolationForest
import joblib

def IsAnomalous(packet_size):
    model_path = "src/ai/anomaly_detection_model.pkl"
    model = joblib.load(model_path)
    prediction = model.predict([[packet_size]])
    return prediction == -1
