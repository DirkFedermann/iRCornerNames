from PyQt5.QtCore import *
from PyQt5.QtWidgets import *
import sys
import threading
import irsdk	# iRacing SDK
import time		# for sleep function
import codecs	# for the german Umlaute (Ä,Ö,Ü) to be displayed correct


class Communicate(QObject):
	signal = pyqtSignal(int, str)

class My_Gui(QWidget):
	def __init__(self):
		super().__init__()

		self.comm = Communicate()
		self.comm.signal.connect(self.append_data)
		self.initUI()

	def initUI(self):

		btn_count = QPushButton('Connect')
		btn_count.clicked.connect(self.start_counting)
		self.te = QLabel()
		self.te.setStyleSheet('font-size:40px;')

		vbox = QVBoxLayout()
		vbox.addWidget(btn_count)
		vbox.addWidget(self.te)

		self.setLayout(vbox)
		self.setWindowTitle('iRcornerNames')
		self.setGeometry(400, 400, 400, 400)
		self.show()

	def count(self, comm):
		ir = irsdk.IRSDK()
		IR_IsActive = ir.startup()
		if IR_IsActive == True:
			comm.signal.emit(0,'Connected')
			f = codecs.open('nordschleife.txt', 'r', 'utf-8')
			corner = f.readline()
			cBegin,cEnd,cName = corner.split(',')
			cBegin = int(cBegin)
			cEnd = int(cEnd)
			while 1 == 1:
				"""
				rpm = str(round(ir['RPM']));
				comm.signal.emit(0,rpm)
				time.sleep(.1)
				"""
				lapDist = round(ir['LapDist'],0)
	
				
				if lapDist > cBegin and lapDist < cEnd:
					comm.signal.emit(0,cName)
					#print(cName)
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
					#print(cName)
					
				time.sleep(.5)
		else:
			comm.signal.emit(0,'cant connect :(')

	def start_counting(self):
		my_Thread = threading.Thread(target=self.count, args=(self.comm,))
		my_Thread.start()

	def append_data(self, num, data):
		self.te.setText(str(num) + " " + data)

if __name__ == '__main__':
	app = QApplication(sys.argv)
	my_gui = My_Gui()
	sys.exit(app.exec_())