namespace go api

service AddService{
    i32 add(1:i32 num1, 2:i32 num2),
    i32 addTimeout(1:i32 num1, 2:i32 num2, 3:i32 timeout),
}