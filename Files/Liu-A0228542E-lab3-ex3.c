#include<string.h>
#include<stdio.h>

int main(){
	__asm__("pushf\norl $0x40000,(%rsp)\npopf");
	char buf[10];
	int *x;
	x = (int *) &buf[1];

	*x = 17;
	printf("*x = %d\n",*x);
	return 0;
}
