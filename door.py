"""
This is deprecated, try daemon.py and systemd
"""

import RPi.GPIO as GPIO
import time

# base command
def start():
  GPIO.setmode(GPIO.BCM)
  GPIO.setwarnings(False)

def clean():
  GPIO.cleanup()
  print 'clean'

# door command
# GPIO pins: https://www.raspberrypi.org/documentation/usage/gpio/images/a-and-b-gpio-numbers.png
def up():
  port=7
  GPIO.setup(port, GPIO.OUT)
  GPIO.output(port, 1)
  time.sleep(0.5)
  GPIO.output(port, 0)
  print 'up'

def stop():
  port=25
  GPIO.setup(port, GPIO.OUT)
  GPIO.output(port, 1)
  time.sleep(0.5)
  GPIO.output(port, 0)
  print 'stop'

def down():
  port=8
  GPIO.setup(port, GPIO.OUT)
  GPIO.output(port, 1)
  time.sleep(0.5)
  GPIO.output(port, 0)
  print 'down'

# input
def input():
  x = raw_input('>')
  if x == 'up':
    up()
  elif x == 'stop':
    stop()
  elif x == 'down':
    down()
  elif x == 'clean':
    clean()

start()
while (True):
  input()
