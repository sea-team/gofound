```

        if(low<i)
            quickSort(arr,low,i-1);//对左端子数组递归
        if(i<high)
            quickSort(arr,j+1,high);//对右端子数组递归
            
        if(high<i) //注意，下标值	 
             quickSort(arr,high,i-1);//对左端子数组递归  
         if(i<low)  //注意，下标值
             quickSort(arr,i+1,low);
```

```java
public static  void quickSort(int[] arr,int high,int low){   
         int i,j,temp;  
         i=high;//高端下标  
         j=low;//低端下标  
         temp=arr[i];//取第一个元素为标准元素。  
           
         while(i<j){//递归出口是 low>=high  
               while(i<j&&temp>arr[j])//后端比temp小，符合降序，不管它，low下标前移
            	   j--;//while完后指比temp大的那个
               if(i<j){
            	   arr[i]=arr[j];
            	   i++;
               }
               while(i<j&&temp<arr[i])
            	   i++;
               if(i<j){
            	   arr[j]=arr[i];
            	   j--;
               }
         }//while完，即第一盘排序  
         arr[i]=temp;//把temp值放到它该在的位置。  
      
         if(high<i) //注意，下标值	 
        	 quickSort(arr,high,i-1);//对左端子数组递归  
         if(i<low)  //注意，下标值
             quickSort(arr,i+1,low);//对右端子数组递归  ；对比上面例子，其实此时i和j是同一下标!!!!!!!!!!!!!
  
     }  

```