import matplotlib.pyplot as plt
import json

# Load the data from 'data.json'
with open('data.json', 'r') as file:
    data = json.load(file)

# Extract 'workers' and 'time' into separate lists
workers = [item['workers'] for item in data]
time = [item['time'] for item in data]

# Create a plot
plt.figure(figsize=(10, 6))
plt.plot(workers, time, marker='o', linestyle='-')
plt.title('Time Taken by UffizziClusters with Varying Workers')
plt.xlabel('Number of Workers')
plt.ylabel('Time')
plt.grid(True)

# Save the plot as an image file
plt.savefig('performance_graph.png')