from sklearn.ensemble import IsolationForest
import joblib
import numpy as np

# Example data - replace with actual feature data
X = np.array([[1500], [1600], [1700], [2000], [100], [200], [300]])

# Train Isolation Forest model
model = IsolationForest(contamination=0.1)
model.fit(X)

# Save the model to a file
joblib.dump(model, 'src/ai/anomaly_detection_model.pkl')
