#!/usr/bin/python

"""
This script reads from a Unix socket at `SOCKET_PATH` and controls the
Rpi GPIO ports.

The log is recorded at `LOG_PATH`.

testing:

    socat - UNIX-CONNECT:/tmp/doorctl
"""

import os
import socket
import sys
import logging

SOCKET_PATH = '/tmp/doorctl'
LOG_PATH = 'doorlog'


class Doorctl:
    def __init__(self):
        """
        setup logger and the Unix domain socket
        """
        logging.basicConfig(filename=LOG_PATH,
                            format='%(asctime)s %(message)s',
                            level=logging.INFO)

        try:
            if os.path.exists(SOCKET_PATH):
                os.unlink(SOCKET_PATH)
        except:
            print('Some errors happend ._.?')
            sys.exit(1)

        self.sock = socket.socket(socket.AF_UNIX)
        self.sock.bind(SOCKET_PATH)
        self.sock.listen(1)

        self.support_commands = ['stop', 'up', 'down']

    def read(self):
        """
        read from the socket
        """
        while True:
            connection, address = self.sock.accept()
            while True:
                data = connection.recv(1024)
                if data:
                    yield data
                else:
                    break

    def run(self, cmd):
        if cmd == 'stop':
            self.stop()
        elif cmd == 'up':
            self.up()
        elif cmd == 'down':
            self.down()

    def stop(self):
        logging.info('Door STOP')
        pass

    def up(self):
        logging.info('Door UP')
        pass

    def down(self):
        logging.info('Door DOWN')
        pass


if __name__ == '__main__':
    doorctl = Doorctl()
    for data in doorctl.read():
        cmd = data.decode().strip()

        if cmd in doorctl.support_commands:
            doorctl.run(cmd)

        else:
            logging.info('Error command from socket: \'{}\''.format(cmd))