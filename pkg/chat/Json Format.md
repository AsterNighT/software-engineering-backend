# Json Format



## Doctor



### Client To Server：



##### Hello

- Type 0
- Role (两种：doctor、patient)
- SenderID

##### MsgFromClient

- Type 1
- SenderID
- ReceiverID
- Content
- Time

##### CloseChat

- Type 2
- PatientID
- DoctorID

##### RequireMedicalRecord

- Type 3
- PatientID

##### RequirePrescription

- Type 4
- PatientID

##### RequireQuestions

- Type 5
- DoctorID





### Server To Client:



##### NewPatient

- Type 6
- PatientID 
- Name

##### MsgFromServer

- Type 7
- SenderID
- Content
- Time

##### SendMedicalRecord

- Type 8
- PatientID
- URL

##### SendPrescription

- Type 9
- PatientID
- URL

##### SendQuestions

- Type 10
- Questions[]







## Patient



### Client To Server：

##### Hello

- Type
- Role (两种：doctor、patient)
- MyID

##### MsgFromClient

- Type
- SenderID
- ReceiverID
- Content
- Time



### Server To Client:

##### MsgFromServer

- Type
- SenderID
- Content
- Time

##### SendMedicalRecord

- Type
- PatientID
- URL

##### SendPrescription

- Type
- PatientID
- URL

##### ChatTerminated

- Type





# Problems

##### ActionListner

##### go的参数

##### 恢复旧信息的机制

点击发起会话，后端查询（senderID，receiverID）是否出现过，若出现过则查询数据库向sender前端推送旧消息，为该（senderID，receiverID）分配相同chatID，调用CreateNewChat与服务器建立连接

##### 数据库的存储、查询

##### 发起新会话的流程

##### 一方发起会话，如何通知另一方接收会话，另一方不在线问题

另一方上线时扫描pending message列表并发送给该用户

##### json





##### 服务端维护list：

ChatID DoctorID  PatientID 病例URL



