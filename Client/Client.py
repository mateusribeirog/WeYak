from SocketHandler import SocketHandler
from gui import GUI
import threading
from queue import Queue
import time

conn = threading.local()
msg = threading.local()
e = threading.Event()

sending_buffer = Queue()
receiving_buffer = Queue()

HOST = "127.0.0.1"
PORT = 8888
app = GUI(sending_buffer,receiving_buffer) 

def Listen(socket,buffer):
    while not e.is_set():
        buff = bytearray(1024)
        temp = socket.receive_message(buff)
        if temp>0:
            print(f"testando o listener: {buff.decode()}")
            buffer.put_nowait(buff[:temp].decode())
            print(f"received, {buff[:temp].decode()}")
        else:
            time.sleep(0.01)
def send(socket,buffer):
    while not e.is_set():
        if not buffer.empty():
            print("not empty")
            msg.value = buffer.get_nowait()
            conn.value = socket.send_message(msg.value.encode('utf-8'))
            print(msg.value)
            if conn.value == 0: 
                print(f"Sent: {msg.value}")
            else:
                print("error at sending")
        else:
            time.sleep(0.01)

def startup_sockets(HOST,PORT):
    s = SocketHandler(HOST,PORT)
    conn.value = s.handle_connection()
    if conn.value == 0:
        print("Sucessfully connected")
        s.send_message(app.Username.encode())
        s.send_message(app.room.encode())
        threading.Thread(target=Listen, args=(s,receiving_buffer),daemon=True).start()
        threading.Thread(target=send, args=(s,sending_buffer),daemon=True).start()
    else:
        print("Connection Failed")
        s.disconnect()
        startup_sockets(HOST, PORT)


def username():
    while not e.is_set():
        if app.Username != "" and app.room != "":
            threading.Thread(target=startup_sockets,args=(HOST,PORT),daemon=True).start()
            return
        else:
            time.sleep(0.01)



def exiting():
    print("Exiting Program")
    e.set()
    app.window.destroy()

threading.Thread(target=username,daemon=True).start()
    
app.window.protocol("WM_DELETE_WINDOW",exiting)




app.run()