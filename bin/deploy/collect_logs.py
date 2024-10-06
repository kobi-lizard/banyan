
import os
import pprint
import matplotlib.pyplot as plt

path = "bin/deploy/logs"
algorithms = ["banyan", "icc", "hotstuff", "streamlet"]
symbols = ['-o', '--+', '-*', '-s']

dictionary = {}

for alg in algorithms:
    dictionary[alg] = []
    files= os.listdir(path)
    files = sorted(files, key=lambda x: os.path.getmtime(os.path.join(path, x)))

    for file in files:
        if os.path.isdir(file) or not file.endswith(".txt") or not file.startswith(alg):
            continue
        s = []
        throughputs = []
        latencies = []
        round_times = []
        totalThroughput = 0.0
        totalLatency = 0.0
        totalRoundTime = 0.0
        fileNo = 0
        roundsNo = 0
        
        f = open(path+"/"+file)
        for line in iter(f):
            if "Throughput" in line:
                throughput = line.strip().split(" ")[-2]
                throughputs.append(throughput)
                totalThroughput += float(throughput)
            if "latency" in line:
                latency = line.strip().split(" ")[-2]
                latencies.append(latency)
                totalLatency += float(latency)
                fileNo += 1
            if "round" in line:
                round_time = line.strip().split(" ")[-2]
                round_times.append(round_time)
                totalRoundTime += float(round_time)
                roundsNo +=1
        f.close()
        print(alg, file.split("_")[1], "Total Throughput:", totalThroughput, "Avg. Latency:", totalLatency / fileNo, "Avg. Round Time:", totalRoundTime / roundsNo)
        dictionary.get(alg).append(tuple((totalThroughput, totalLatency / fileNo)))
# with open(path+'/'+str(fileNo)+".log", 'w') as w_f:
#     w_f.write("Experiment with " + str(fileNo) +" clients:" + "\n")
#     for i in range(fileNo):
#         w_f.write("["+str(i)+"]"+" throughput: "+str(throughputs[i])+", latency: "+str(latencies[i])+"\n")

#     w_f.write("Total throughput: "+str(totalThroughput))
#     w_f.write("\nAverage latency: "+str(totalLatency/fileNo))

print(dictionary)
with open("bin/deploy/logs/new.data", "w") as file:
    pprint.pprint(dictionary, stream=file)
    

data = []
for key, value in dictionary.items():
    data.append(tuple((key, value, symbols.pop())))


def do_plot():
    f = plt.figure(1, figsize=(7,5));
    plt.clf()
    ax = f.add_subplot(1, 1, 1)
    for name, entries, style in data:
        throughput = []
        latency = []
        for t, l in entries:
            # batch.append(N*ToverN)
            # throughput.append(ToverN*(N-t) / latency)
            throughput.append(t)
            latency.append(l)
        ax.plot(throughput, latency, style, label='%s' % name)
    #ax.set_xscale("log")
#     ax.set_yscale("log")
    # plt.ylim([0, 50])
    #plt.xlim([10**3.8, 10**6.4])
    plt.legend(loc='upper left')
    # plt.ylabel('Throughput (Tx per second) in log scale')
    plt.ylabel('Latency (ms)')
    plt.xlabel('Throughput (KTx/s)')
    # plt.xlabel('Requests (Tx) in log scale')
    plt.tight_layout()
    # plt.show()
    plt.savefig('bin/deploy/logs/happy-path.png', format='png', dpi=400)

if __name__ == '__main__':
    do_plot()
