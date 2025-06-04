import tkinter as tk
from ctypes import windll
from queue import Queue

class GUI:
    def __init__(self,sending_buff, receiving_buff):
        
        self.PAGE_NAMES = ["Games","Tech","General","Movies"]
        self.Username = ""
        self.room = ""
        self.sending_buffer = sending_buff
        self.receiving_buff = receiving_buff

        windll.shcore.SetProcessDpiAwareness(1)
        #Configuração inicial da pagina
        self.window = tk.Tk()
        self.window.geometry("800x800")
        self.window.resizable(False,False)
        self.window.title("WeYak!")
        self.icon = tk.PhotoImage(file="logo.png")
        self.window.iconphoto(False,self.icon)
        #Configuração dos botoes
        self.topnotch = tk.Frame(master=self.window,borderwidth=1)
        self.topnotch.grid(row=0, column=0,sticky="nsew")
        self.window.grid_columnconfigure(0,weight=1)

        self.label = tk.Label(self.topnotch, text="WeYak!", font=("Arial", 18),bg="#35c257")
        self.label.pack(fill="both",expand=True)
        self.buttons = []

        for i in range(4):
            Button = tk.Button(text=self.PAGE_NAMES[i],width=1,height=1,relief=tk.RAISED,font=("Arial",14),command=eval(f"self.set_room_{self.PAGE_NAMES[i]}"))
            Button.grid(row=0,column=i+1,sticky="nsew")
            self.window.grid_columnconfigure(i+1,weight=1)
            self.buttons.append(Button)


        self.window.grid_rowconfigure(1,weight=6)
        self.window.grid_rowconfigure(2,weight=1)

        #Configuração do Chat

        self.Text = tk.Listbox(master=self.window, width=1,height=1,font=("Arial",10))
        self.Text.grid(row=1,column=0,columnspan=3 ,sticky="nsew")
        self.scroll = tk.Scrollbar(self.window, command=self.Text.yview)
        self.Text.configure(yscrollcommand=self.scroll.set)
        self.scroll.grid(row=1, column=3, sticky="ns")

        #input:
        self.input_frame = tk.Frame(self.window)
        self.input_frame.grid(row=2,column=0, columnspan=3, sticky="nsew")

        #label do input
        self.input_label = tk.Label(self.input_frame,text="Write here:", font=("Arial",12))
        self.input_label.pack(side="top")

        #Configuração do input
        self.input = tk.Entry(self.input_frame,width=1,state="disabled")
        self.input.pack(fill=tk.X)
        self.input.bind("<Return>",self.sendMsg)
        self.window.grid_rowconfigure(2,weight=1)

        #RightSideFrame:
        self.rightFrame = tk.Frame(self.window)
        self.rightFrame.grid(column=3,row=1,rowspan=2,columnspan=2,sticky="nsew")

        #Current Room:
        self.current_Room = tk.Label(self.rightFrame,text="Current Room:\n",font=("Arial",12))
        self.current_Room.pack(pady=(20,100))


        #UserName:

        self.Username_frame = tk.Frame(self.rightFrame)
        self.Username_frame.pack()

        self.Username_label = tk.Label(self.Username_frame, text="Username:",font=("Arial",12))
        self.Username_label.pack()
        self.Username_input = tk.Entry(self.Username_frame)
        self.Username_input.pack()

        self.window.after(100,self.update_TextBox)

        self.Username_input.bind("<Return>",self.setUsername)


        #Server Status:
        self.server_status = tk.Label(self.rightFrame,text="teste",font=("Arial",14))
        self.server_status.pack(pady=(100,0))

    def set_room_General(self):
        self.room = "General\n"
        self.input.config(state="normal")
    def set_room_Games(self):
        self.room = "Games\n"
        self.input.config(state="normal")
    def set_room_Tech(self):
        self.room = "General\n"
        self.input.config(state="normal")
    def set_room_Movies(self):
        self.room = "General\n"
        self.input.config(state="normal")

    def setUsername(self,event):
        self.Username = self.Username_input.get() + '!'
        self.Username_input.config(state='readonly')

    def sendMsg(self,event):
        msg = "{}\n".format(self.input.get())
        #self.receiving_buff.put_nowait(msg)
        self.sending_buffer.put_nowait(msg)
        self.input.delete(0,'end')
        self.update_TextBox()
    
    def update_TextBox(self):
        while not self.receiving_buff.empty():
            print(f"size queue: {self.receiving_buff.qsize()}")
            buffer = self.receiving_buff.get_nowait()
            buffer = buffer.split(":")
            print(f"o buffer: {buffer}")
            print(f"size queue: {self.receiving_buff.qsize()}")
            try:
                #self.Text.config(state="normal")
                msg = f"{buffer[0]}: {buffer[1]}"
                self.Text.insert('end',msg+'\n')
                print(f"texto da textBox:{msg}")
                self.Text.yview_moveto(1.0)
            except Exception as e:
                print(f"erro no tk: {e}")
        self.window.after(100,self.update_TextBox)


    def check_new_msgs(self):
        if not self.receiving_buff.empty():
            print("nao ta vazio")
            self.update_TextBox()
        self.window.after(1000,self.check_new_msgs)
             

    def run(self):
        self.window.mainloop()