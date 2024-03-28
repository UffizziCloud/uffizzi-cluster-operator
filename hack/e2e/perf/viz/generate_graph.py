import matplotlib.pyplot as plt
import json

# Load data from JSON
with open('data.json', 'r') as file:
    data = json.load(file)

# Extracting workers and time into separate lists
workers = [item['workers'] for item in data]
time = [item['time'] for item in data]

# Plotting
plt.figure(figsize=(10, 6))
plt.plot(workers, time, marker='o')
plt.title('Test Performance')
plt.xlabel('Number of Workers')
plt.ylabel('Time')
plt.grid(True)
plt.savefig('output.png')