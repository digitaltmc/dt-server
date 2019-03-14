package main

func makeMeetingItemResolver(agenda []*MeetingItem) []*meetingItemResolver {
  ret := make([]*meetingItemResolver, len(agenda))
  for i, v := range agenda {
    ret[i] = &meetingItemResolver{v}
  }
  return ret
}


