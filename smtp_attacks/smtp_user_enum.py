#!/usr/bin/env python3
"""Uses the SMTP VRFY command to enumerate usernames (i.e. mailaddresses)
"""
import sys,socket,time

FILENAME = str(sys.argv[1])
HOSTNAME = str(sys.argv[2])
PORT = 25

def receive_data(socket, timeout=15):
    start_time = time.time()
    response = b''
    while time.time() - start_time < timeout:
        try:
            received = socket.recv(1024)
            if not received:
                raise Exception('Connection closed by sender')
            response += received
        except Exception:
            pass
    print(response)
    return response.decode("utf-8").strip()


with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as smtp_socket:
    smtp_socket.connect((HOSTNAME, PORT))
    smtp_socket.settimeout(1)

    smtp_socket.sendall(b"HELO myself")
    response = receive_data(smtp_socket)
    assert "SMTP" in response

    print("start enum")
    for line in open(FILENAME,"r").readlines():
        smtp_socket.sendall(bytes("VRFY %s"%(line.strip()), 'utf-8'))
        response = receive_data(smtp_socket)
        print(response)
        if response.startswith("252 "):
            print("Found username %s"%line.strip())
        else:
            assert response.startswith("550 ")
    