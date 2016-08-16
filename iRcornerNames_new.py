#!python3
# -*- coding: UTF-8 -*-

import irsdk	# iRacing SDK
import time		# for sleep function
import os		# for cls - clear screen
import codecs	# for the german Umlaute (Ä,Ö,Ü) to be displayed correct


# initiate the SDK
ir = irsdk.IRSDK()
IR_IsActive = ir.startup()

# check if iRacing is running
while IR_IsActive != True:
	i = 10
	while i > 0:
		print('iRacing is not running... I try it in ' + i + ' seconds again')
		sleep(1)
		i = i - 1
	IR_IsActive = ir.startup()
	
print('iRacing is running')

f = codecs.open('nordschleife.txt', 'r', 'utf-8')
corner = f.readline()
cBegin,cEnd,cName = corner.split(',')
cBegin = int(cBegin)
cEnd = int(cEnd)

while IR_IsActive != False:
	
	lapDist = round(ir['LapDist'],0)
	
	os.system('cls')
	
	if lapDist > cBegin and lapDist < cEnd:
		print(cName)
	else:
		corner = f.readline()
		cBegin,cEnd,cName = corner.split(',')
		cBegin = int(cBegin)
		cEnd = int(cEnd)
		if cName == 'EOF':
			f.close()
			f = codecs.open('nordschleife.txt', 'r', 'utf-8')
			corner = f.readline()
			cBegin,cEnd,cName = corner.split(',')
			cBegin = int(cBegin)
			cEnd = int(cEnd)
		print(cName)
		
	time.sleep(.5)
	
f.close()