def hello(event, context):
    print(event)

    return event['data']


def world(event, context):
    print(event)

    return event['data']
