# wakeonlan-go

<div align="right">

[![Codacy Badge](https://app.codacy.com/project/badge/Grade/1ad5e14314a44a37a287969eef46bd40)](https://www.codacy.com/gh/DanArmor/wakeonlan-go/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=DanArmor/wakeonlan-go&amp;utm_campaign=Badge_Grade)

</div>

WoL tool

#### Usage 
```wakeonlan-go -m <mac-address>``` - default usage, should work fine in most cases

```wakeonlan-go -m <mac-address> -d 255.255.255.255:9``` - when you need to specify destination

```wakeonlan-go -m <mac-address> -l 192.168.10.100:9``` - when you need to specify your networking interface, from which you're sending WoL packet.

On Windows you can get list of your networking interfaces and their addresses with `ipconfig`

On Linux you can do it with `ip address show`

#### To be done
*   Aliases 
*   Flags
    -   Port as separate flag
