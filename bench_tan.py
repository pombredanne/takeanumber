import time

from takeanumber import Client


c = Client()


def bench():
    start = time.time()

    for i in xrange(1000000):
        c.add('my_queue', 'hello #{}'.format(i))

        if i % 10000 == 0:
            print i

    elapsed = time.time() - start
    print "Elapsed time for 1000000: {}".format(elapsed)
    print c.len('my_queue')


bench()
