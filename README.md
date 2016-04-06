Hermes
========

rsyslog -> doppler -> traffic controller -> clients

# Doppler
Log Receiver -> Diode Shards -> (Single Diode Shard / Individual Go Routine) -> Map to Linked List of Listeners

### Map to Linked List of Listeners
Accessed via 3 different go routines
* Added to via kv-store event (Write go routine)
* Deleted from via ttl (Delete go routine)
* Read from via Log Flow (Read go routine)
