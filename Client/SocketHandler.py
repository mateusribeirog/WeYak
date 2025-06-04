""" import socket

HOST = "127.0.0.1"
PORT = 12000

with socket.socket(socket.AF_INET,socket.SOCK_STREAM) as s:
    s.connect((HOST,PORT))
    
    while True:
        try:
            sender = input("")
            s.sendall(sender.encode('utf-8'))
            echo = s.recv(1024)
            print("echo: ",echo)
        except KeyboardInterrupt:
            print("conexão finalizada pelo cliente")
            break """
import socket

class SocketHandler:
    def __init__(self, HOST, PORT,sock=None):
        self.HOST = HOST
        self.PORT = PORT
        if sock is None:
            self.socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        else:
            self.socket = sock
        print(f"Socket Created at: {HOST,PORT}")
    
    def handle_connection(self):
        try:
            print(f"Trying to connect to: {self.HOST,self.PORT}")
            self.socket.connect((self.HOST, self.PORT))
            return 0
        except OSError as e:
            print(f"Error in connection: {e}")
            return -1
    
    def disconnect(self):
        try:
            self.socket.close()
        except OSError as e:
            print(f"Error in disconnecting from socket: {e}")
            return -1
    
    def send_message(self,buff):
        if len(buff) > 0:
            try:
                self.socket.sendall(buff)
                return 0
            except OSError as e:
                print(f"Error in sending messages: {e}")
                return -1
    
    def receive_message(self,buff):
        try:
            temp = self.socket.recv_into(buff,1024)
            return temp
        except OSError as e:
            print(f"Error at sending messages: {e}")
            return -1
