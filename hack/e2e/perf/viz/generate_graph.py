import json
import matplotlib.pyplot as plt

# Load the data from data.json
with open('data.json', 'r') as file:
    data = json.load(file)

# Extract workers and time data
workers = [item['workers'] for item in data]
time = [item['time'] for item in data]

# Create a graph
plt.figure(figsize=(10, 6))
plt.plot(workers, time, marker='o')
plt.title('Time taken to create UffizziClusters')
plt.xlabel('Number of Workers')
plt.ylabel('Time (s)')
plt.grid(True)
plt.savefig('graph.png')