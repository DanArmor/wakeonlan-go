# wakeonlan-go

<div align="right">

[![Codacy Badge](https://app.codacy.com/project/badge/Grade/1ad5e14314a44a37a287969eef46bd40)](https://www.codacy.com/gh/DanArmor/wakeonlan-go/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=DanArmor/wakeonlan-go&amp;utm_campaign=Badge_Grade)

</div>

WoL tool

#### Usage 
```wakeonlan-go [OPTIONS] MACs...``` - where ```MACs``` can be addresses or aliases

  `-d=`         Destination address(port is optional)
  
  `-l=`         Local address of network interface (port is optional)
  
  `-r=`         If used - action is not performed, but is recorded as alias to given name. If you use alias inside other alias - it will be deep copy,
              so you can delete used alias in the future.
              
  `-s`          Show list of aliases
  
      `--rm=`   Remove alias with a given name

On Windows you can get list of your networking interfaces and their addresses with `ipconfig` to use with `-l`/`-d` flags.

On Linux you can do it with `ip address show`

#### Screenshots

![image](https://user-images.githubusercontent.com/39347109/227588346-c22f7a40-dcb1-4ceb-8147-72dd6b737d93.png)

![image](https://user-images.githubusercontent.com/39347109/227588667-460df642-2320-4f15-a55c-210d1efd9b0c.png)

![image](https://user-images.githubusercontent.com/39347109/227588943-6cf841a5-0232-4751-b33d-4681ec51c96a.png)


