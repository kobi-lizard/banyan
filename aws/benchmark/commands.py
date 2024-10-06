# Copyright(C) Facebook, Inc. and its affiliates.
from os.path import join

from benchmark.utils import PathMaker


class CommandMaker:

    @staticmethod
    def cleanup():
        return (
            f'rm -r .db-* ; rm .*.json ; mkdir -p {PathMaker.results_path()}'
        )

    @staticmethod
    def clean_logs():
        return f'rm -r {PathMaker.logs_path()} ; mkdir -p {PathMaker.logs_path()}'

    @staticmethod
    def compile():
        return 'cargo build --quiet --release'

    @staticmethod
    def generate_key(filename):
        assert isinstance(filename, str)
        return f'./node generate_keys --filename {filename}'

    @staticmethod
    def generate_config_files(bport,client_bport,client_run_port,num_nodes):
        return f'./genconfig --blocksize 100 --delay 100 --base_port {bport} --client_base_port {client_bport} --NumNodes {num_nodes} --target . --client_run_port {client_run_port} --local true'

    @staticmethod
    def run_primary(key,delay,ep,delta,val,tri,batch,freq, debug=False):
        assert isinstance(key, str)
        assert isinstance(debug, bool)
        #v = '-vvv' if debug else '-vv'
        return (f'./node --config {key} --ip ip_file '
                f'--sleep {delay} --vsstype glow --epsilon {ep} --delta {delta} --val {val} --tri {tri} --syncer syncer --batch {batch} --frequency {freq}')
    
    @staticmethod
    def run_syncer(key,delay, debug=False):
        assert isinstance(key, str)
        assert isinstance(debug, bool)
        #v = '-vvv' if debug else '-vv'
        return (f'./node --config {key} --ip ip_file '
                f'--sleep {delay} --vsstype sync --epsilon 10 --delta 5000 --val 100 --tri 10 --syncer syncer --batch 200 --frequency 10')

    @staticmethod
    def unzip_tkeys(fileloc,dir, debug=False):
        return (f'tar -xvzf {fileloc} && cp {dir}/* .')

    @staticmethod
    def run_worker(keys, committee, store, parameters, id, debug=False):
        assert isinstance(keys, str)
        assert isinstance(committee, str)
        assert isinstance(parameters, str)
        assert isinstance(debug, bool)
        v = '-vvv' if debug else '-vv'
        return (f'./node {v} run --keys {keys} --committee {committee} '
                f'--store {store} --parameters {parameters} worker --id {id}')

    @staticmethod
    def run_client(address, size, rate, nodes):
        assert isinstance(address, str)
        assert isinstance(size, int) and size > 0
        assert isinstance(rate, int) and rate >= 0
        assert isinstance(nodes, list)
        assert all(isinstance(x, str) for x in nodes)
        nodes = f'--nodes {" ".join(nodes)}' if nodes else ''
        return f'./benchmark_client {address} --size {size} --rate {rate} {nodes}'

    @staticmethod
    def kill():
        return 'tmux kill-server'

    @staticmethod
    def alias_binaries(origin):
        assert isinstance(origin, str)
        node, client, genconfig = join(origin, 'node'), join(origin, 'benchmark_client'), join(origin,'genconfig')
        return f'rm node ; rm benchmark_client ; rm genconfig ; ln -s {node} . ; ln -s {client} . ; ln -s {genconfig} .'
