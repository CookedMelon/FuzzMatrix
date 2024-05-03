# -*- coding: utf-8 -*-
import tornado
import tornado.websocket
import paramiko
import threading
import time

# ���÷�������Ϣ
HOSTS = ""
PORT = 22
USERNAME = ""
PASSWORD = ""


class MyThread(threading.Thread):
    def __init__(self, id, chan):
        threading.Thread.__init__(self)
        self.chan = chan


    def run(self):
        # �����߳��д����¼�ѭ��
        tornado.ioloop.IOLoop().make_current()
        self.ioloop = tornado.ioloop.IOLoop.current()

        while not self.chan.chan.exit_status_ready():
            time.sleep(0.1)
            try:
                data = self.chan.chan.recv(1024)
                self.chan.write_message(data)
            except Exception as ex:
                print(str(ex))
        self.chan.sshclient.close()
        return False



class webSSHServer(tornado.websocket.WebSocketHandler):

    def open(self):
        self.sshclient = paramiko.SSHClient()
        self.sshclient.load_system_host_keys()
        self.sshclient.set_missing_host_key_policy(paramiko.AutoAddPolicy())
        self.sshclient.connect(HOSTS, PORT, USERNAME, PASSWORD)
        self.chan = self.sshclient.invoke_shell(term='xterm')
        self.chan.settimeout(0)
        t1 = MyThread(999, self)
        t1.daemon = True
        t1.start()

    def on_message(self, message):
        try:
            self.chan.send(message)
        except Exception as ex:
            print(str(ex))

    def on_close(self):
        self.sshclient.close()

    def check_origin(self, origin):
        # ����������
        return True


if __name__ == '__main__':
    # ����·��
    app = tornado.web.Application([
        (r"/terminals", webSSHServer),
    ],
        debug=True
    )

    # ����������
    http_server = tornado.httpserver.HTTPServer(app)
    http_server.listen(3000)
    tornado.ioloop.IOLoop.current().start()
