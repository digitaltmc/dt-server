#!/bin/bash

# WARNING: drop all indexes except the _id one.
db.meeting.dropIndexes()

db.meeting.createIndex({date: 1, roleName: 1}, { unique: true })

db.meeting.getIndexes()
