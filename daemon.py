#!/usr/bin/python

"""
This script reads from a Unix socket at `SOCKET_PATH` and controls the
Rpi GPIO ports.

The log is recorded at `LOG_PATH`.

testing:

    socat - UNIX-CONNECT:/tmp/doorctl

GPIO pins:

https://www.raspberrypi.org/documentation/usage/gpio/images/a-and-b-gpio-numbers.png

GPIO signal:

    |port |roll-up door action|
    |-------------------------|
    |7    | UP                |
    |25   | STOP              |
    |8    | DOWN              |
"""

import RPi.GPIO as GPIO

import logging
import os
import socket
import sys
import time

SOCKET_PATH = '/tmp/doorctl'
LOG_PATH = '/tmp/doorlog'


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

        self.support_commands = ['stop', 'up', 'down', 'clean']

        # GPIO setup
        GPIO.setmode(GPIO.BCM)
        GPIO.setwarnings(False)

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
        elif cmd == 'clean':
            self.clean()

    def signal(self, port):
        """
        raise GPIO signal for 0.5 second on the specific port
        """
        logging.debug('signal on port {}'.format(port))
        GPIO.setup(port, GPIO.OUT)
        GPIO.output(port, 1)
        time.sleep(0.5)
        GPIO.output(port, 0)

    def stop(self):
        logging.info('Door STOP')
        self.signal(port=25)

    def up(self):
        logging.info('Door UP')
        self.signal(port=7)

    def down(self):
        logging.info('Door DOWN')
        self.signal(port=8)

    def clean(self):
        GPIO.cleanup()


if __name__ == '__main__':
    doorctl = Doorctl()
    for data in doorctl.read():
        cmd = data.decode().strip()

        if cmd in doorctl.support_commands:
            doorctl.run(cmd)

        else:
            logging.info('Error command from socket: \'{}\''.format(cmd))
