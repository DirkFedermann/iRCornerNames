from PyQt5.QtCore import *
from PyQt5.QtWidgets import *
from PyQt5.QtGui import *
import sys
import threading
import irsdk	# iRacing SDK
import time		# for sleep function
import codecs	# for the german Umlaute (Ä,Ö,Ü) to be displayed correct

# read the settings.ini
config = {}
exec(open("settings.conf").read(), config)


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
		self.te.setStyleSheet('font-size:' + str(config['fontSize']) + 'px;')
		
		if config['debug'] == 1:
			self.debug = QLabel()
			self.debug.setStyleSheet('font-size:' + str(config['fontSize']) + 'px;')

		vbox = QVBoxLayout()
		vbox.addWidget(btn_count)
		vbox.addWidget(self.te)
		if config['debug'] == 1:
			vbox.addWidget(self.debug)

		self.setLayout(vbox)
		self.setWindowTitle('iRcornerNames')
		self.setGeometry(config['window_X'], config['window_Y'],config['window_Width'], config['window_Height'])
		self.show()

	def count(self, comm):
		ir = irsdk.IRSDK()
		IR_IsActive = ir.startup()
		if IR_IsActive == True:
			comm.signal.emit(0,'Connected')
			f = codecs.open(config['track'], 'r', 'utf-8')
			corner = f.readline()
			cBegin,cEnd,cName = corner.split(',')
			cBegin = int(cBegin)
			cEnd = int(cEnd)
			while IR_IsActive == True:
				lapDist = round(ir['LapDist'],0)
				
				if config['debug'] == 1:
					comm.signal.emit(1,str(lapDist))
				
				if lapDist > cBegin and lapDist < cEnd:
					comm.signal.emit(0,cName)
				else:
					corner = f.readline()
					cBegin,cEnd,cName = corner.split(',')
					cBegin = int(cBegin)
					cEnd = int(cEnd)
					if cName == 'EOF':
						f.close()
						f = codecs.open(config['track'], 'r', 'utf-8')
						corner = f.readline()
						cBegin,cEnd,cName = corner.split(',')
						cBegin = int(cBegin)
						cEnd = int(cEnd)
					
				time.sleep(config['update_Time'])
		else:
			comm.signal.emit(0,'cant connect :(')

	def start_counting(self):
		my_Thread = threading.Thread(target=self.count, args=(self.comm,))
		my_Thread.start()

	def append_data(self, num, data):
		self.te.setText(data)
		if config['debug'] == 1 and num == 1:
			self.debug.setText(data)

			
if __name__ == '__main__':
	app = QApplication(sys.argv)
	my_gui = My_Gui()
	sys.exit(app.exec_())
