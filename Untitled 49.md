```c
#include "kernel/types.h"
#include "kernel/stat.h"
#include "user/user.h"

void childPipe(int fds[2]);
int main() {
    int fds1[2];
    pipe(fds1);
    if(fork()==0){
        childPipe(fds1);
    }else{  
        close(fds1[0]);
     for(int i = 2; i <= 35; i++){
        
        
            write(fds1[1],&i,sizeof(int));

        }
        close(fds1[1]);
        wait(0);
    


    }
    
    exit(0);
}
// 每次从fds中读取素数，并把读取后的传递给下一次·
void childPipe(int fds[2]){
   int prime;
    int next;
    // 不先关闭的话 读出来会有很多乱码 
    close(fds[1]);
    if (read(fds[0], &prime, sizeof(int))) {
        printf("prime %d\n", prime);
        int fds2[2];
        pipe(fds2);
    
        if (fork() == 0) {
            childPipe(fds2);
        } else {
            close(fds2[0]);
            while (read(fds[0], &next, sizeof(int))) {
                // 下一次需要写入的素数
                if (next % prime != 0)
                    write(fds2[1], &next, sizeof(int));
            }
            
            close(fds2[1]);
            wait(0);
        }
        exit(0); // 移到条件语句内
    }
    exit(0);
    }
    



// // 从上一个管道读入数据 0读入，1写入
// void pipeFork(int fds[2],int i,int recive) {
//     if(i == 31) {
//         read(fds[0],&recive,4);
//         return;
//     }
//     if(isPrime(i)){
//         if(fork()==0){
//             int fds1[2];
//             pipe(fds1);
//             // 向下一个进程的管道写入
//             write(fds1[0],&recive,4);
//             pipeFork(fds1,i+1, recive);
//         }else{
//             read(fds1[0],&recive,4);
//         }
//     }


// }

```





```c
#include "kernel/types.h"
#include "kernel/stat.h"
#include "user/user.h"
int main(int arg, char* argv[]) {
 
    int fds[2];
    pipe(fds);
    // pipe两端 一端收，一端发
    if(fork() == 0){
        int id = getpid();
        read(fds[0],&id,1);
        fprintf(1,"%d: received ping\n", getpid());

        write(fds[1],&id,1);
        
        // read(0,buf,sizeof(buf));

    }else {
        int id = getpid();
        write(fds[1],&id,1);
        wait(0);
        read(fds[0],&id,1);
        fprintf(1,"%d: received pong\n", getpid());
    }
    close(fds[1]);
    close(fds[0]);

    exit(0);

}
```





```
#include "kernel/types.h"
#include "kernel/stat.h"
#include "user/user.h"
int 
main(int argc, char * argv[]) 
{
    if (argc != 2)
        {
                fprintf(2, "Usage: sleep seconds\n");
                exit(1);
        }
        int time = atoi(argv[1]);
        sleep(time);

        exit(0);
}
```





```
#include "kernel/param.h"
#include "kernel/types.h"
#include "kernel/stat.h"
#include "user/user.h"

int
main(int argc, char *argv[])
{
  int i;
  char *nargv[MAXARG];

  if(argc < 3 || (argv[1][0] < '0' || argv[1][0] > '9')){
    fprintf(2, "Usage: %s mask command\n", argv[0]);
    exit(1);
  }

  if (trace(atoi(argv[1])) < 0) {
    fprintf(2, "%s: trace failed\n", argv[0]);
    exit(1);
  }
  
  for(i = 2; i < argc && i < MAXARG; i++){
    nargv[i-2] = argv[i];
  }
  exec(nargv[0], nargv);
  exit(0);
}




#include "types.h"
#include "param.h"
#include "memlayout.h"
#include "riscv.h"
#include "spinlock.h"
#include "proc.h"
#include "syscall.h"
#include "defs.h"

// Fetch the uint64 at addr from the current process.
int
fetchaddr(uint64 addr, uint64 *ip)
{
  struct proc *p = myproc();
  if(addr >= p->sz || addr+sizeof(uint64) > p->sz)
    return -1;
  if(copyin(p->pagetable, (char *)ip, addr, sizeof(*ip)) != 0)
    return -1;
  return 0;
}

// Fetch the nul-terminated string at addr from the current process.
// Returns length of string, not including nul, or -1 for error.
int
fetchstr(uint64 addr, char *buf, int max)
{
  struct proc *p = myproc();
  int err = copyinstr(p->pagetable, buf, addr, max);
  if(err < 0)
    return err;
  return strlen(buf);
}

static uint64
argraw(int n)
{
  struct proc *p = myproc();
  switch (n) {
  case 0:
    return p->trapframe->a0;
  case 1:
    return p->trapframe->a1;
  case 2:
    return p->trapframe->a2;
  case 3:
    return p->trapframe->a3;
  case 4:
    return p->trapframe->a4;
  case 5:
    return p->trapframe->a5;
  }
  panic("argraw");
  return -1;
}

// Fetch the nth 32-bit system call argument.
int
argint(int n, int *ip)
{
  *ip = argraw(n);
  return 0;
}

// Retrieve an argument as a pointer.
// Doesn't check for legality, since
// copyin/copyout will do that.
int
argaddr(int n, uint64 *ip)
{
  *ip = argraw(n);
  return 0;
}

// Fetch the nth word-sized system call argument as a null-terminated string.
// Copies into buf, at most max.
// Returns string length if OK (including nul), -1 if error.
int
argstr(int n, char *buf, int max)
{
  uint64 addr;
  if(argaddr(n, &addr) < 0)
    return -1;
  return fetchstr(addr, buf, max);
}

extern uint64 sys_chdir(void);
extern uint64 sys_close(void);
extern uint64 sys_dup(void);
extern uint64 sys_exec(void);
extern uint64 sys_exit(void);
extern uint64 sys_fork(void);
extern uint64 sys_fstat(void);
extern uint64 sys_getpid(void);
extern uint64 sys_kill(void);
extern uint64 sys_link(void);
extern uint64 sys_mkdir(void);
extern uint64 sys_mknod(void);
extern uint64 sys_open(void);
extern uint64 sys_pipe(void);
extern uint64 sys_read(void);
extern uint64 sys_sbrk(void);
extern uint64 sys_sleep(void);
extern uint64 sys_unlink(void);
extern uint64 sys_wait(void);
extern uint64 sys_write(void);
extern uint64 sys_uptime(void);
extern uint64 sys_trace(void);
extern uint64 sys_info(void);

static uint64 (*syscalls[])(void) = {
[SYS_fork]    sys_fork,
[SYS_exit]    sys_exit,
[SYS_wait]    sys_wait,
[SYS_pipe]    sys_pipe,
[SYS_read]    sys_read,
[SYS_kill]    sys_kill,
[SYS_exec]    sys_exec,
[SYS_fstat]   sys_fstat,
[SYS_chdir]   sys_chdir,
[SYS_dup]     sys_dup,
[SYS_getpid]  sys_getpid,
[SYS_sbrk]    sys_sbrk,
[SYS_sleep]   sys_sleep,
[SYS_uptime]  sys_uptime,
[SYS_open]    sys_open,
[SYS_write]   sys_write,
[SYS_mknod]   sys_mknod,
[SYS_unlink]  sys_unlink,
[SYS_link]    sys_link,
[SYS_mkdir]   sys_mkdir,
[SYS_close]   sys_close,
[SYS_trace]   sys_trace, // NEW
[SYS_sysinfo]   sys_info, // NEW
};

static char* syscall_names[] = {
  [SYS_fork]    "fork",
  [SYS_exit]    "exit",
  [SYS_wait]    "wait",
  [SYS_pipe]    "pipe",
  [SYS_read]    "read",
  [SYS_kill]    "kill",
  [SYS_exec]    "exec",
  [SYS_fstat]   "fstat",
  [SYS_chdir]   "chdir",
  [SYS_dup]     "dup",
  [SYS_getpid]  "getpid",
  [SYS_sbrk]    "sbrk",
  [SYS_sleep]   "sleep",
  [SYS_uptime]  "uptime",
  [SYS_open]    "open",
  [SYS_write]   "write",
  [SYS_mknod]   "mknod",
  [SYS_unlink]  "unlink",
  [SYS_link]    "link",
  [SYS_mkdir]   "mkdir",
  [SYS_close]   "close",
  [SYS_trace]   "trace",
  [SYS_sysinfo] "sys_info",
};
void
syscall(void)
{
  int num;
  struct proc *p = myproc();
  char* syscall_name; 
  num = p->trapframe->a7;
  
  if(num > 0 && num < NELEM(syscalls) && syscalls[num]) {
    p->trapframe->a0 = syscalls[num]();
    if((p->trace_mask &(1 << num))!=0){
       syscall_name = syscall_names[num];                                   
        printf("%d: syscall %s -> %d\n", p->pid, syscall_name, p->trapframe->a0);
    }
  } else {
    printf("%d %s: unknown sys call %d\n",
            p->pid, p->name, num);
    p->trapframe->a0 = -1;
  }
}


#include "types.h"
#include "riscv.h"
#include "defs.h"
#include "date.h"
#include "param.h"
#include "memlayout.h"
#include "spinlock.h"
#include "proc.h"
#include "sysinfo.h"

struct sysinfo;

uint64
sys_exit(void)
{
  int n;
  if(argint(0, &n) < 0)
    return -1;
  exit(n);
  return 0;  // not reached
}

uint64
sys_getpid(void)
{
  return myproc()->pid;
}

uint64
sys_fork(void)
{
  return fork();
}

uint64
sys_wait(void)
{
  uint64 p;
  if(argaddr(0, &p) < 0)
    return -1;
  return wait(p);
}

uint64
sys_sbrk(void)
{
  int addr;
  int n;

  if(argint(0, &n) < 0)
    return -1;
  addr = myproc()->sz;
  if(growproc(n) < 0)
    return -1;
  return addr;
}

uint64
sys_sleep(void)
{
  int n;
  uint ticks0;

  if(argint(0, &n) < 0)
    return -1;
  acquire(&tickslock);
  ticks0 = ticks;
  while(ticks - ticks0 < n){
    if(myproc()->killed){
      release(&tickslock);
      return -1;
    }
    sleep(&ticks, &tickslock);
  }
  release(&tickslock);
  return 0;
}

uint64
sys_kill(void)
{
  int pid;

  if(argint(0, &pid) < 0)
    return -1;
  return kill(pid);
}

// return how many clock tick interrupts have occurred
// since start.
uint64
sys_uptime(void)
{
  uint xticks;

  acquire(&tickslock);
  xticks = ticks;
  release(&tickslock);
  return xticks;
}
uint64
sys_trace(void)
{
  int n;
  if(argint(0,&n) < 0)
    return -1;
  myproc()->trace_mask = n;
  return 0;
}

uint64
sys_info(void)
{
  
  struct sysinfo p;
  uint64 st; // user pointer to struct stat
  p.freemem = free_mem();
  p.nproc = freeProc();
  if(argaddr(0, &st) < 0)
    return -1;

  if(copyout(myproc()->pagetable, st, (char *)&p, sizeof(p)) < 0)
      return -1;
  else return 0;

}




```



```
#include "types.h"
#include "riscv.h"
#include "param.h"
#include "defs.h"
#include "memlayout.h"
#include "spinlock.h"
#include "proc.h"

uint64
sys_exit(void)
{
  int n;
  argint(0, &n);
  exit(n);
  return 0;  // not reached
}

uint64
sys_getpid(void)
{
  return myproc()->pid;
}

uint64
sys_fork(void)
{
  return fork();
}

uint64
sys_wait(void)
{
  uint64 p;
  argaddr(0, &p);
  return wait(p);
}

uint64
sys_sbrk(void)
{
  uint64 addr;
  int n;

  argint(0, &n);
  addr = myproc()->sz;
  if(growproc(n) < 0)
    return -1;
  return addr;
}

uint64
sys_sleep(void)
{
  int n;
  uint ticks0;
  argint(0, &n);
  acquire(&tickslock);
  ticks0 = ticks;
  while(ticks - ticks0 < n){
    if(killed(myproc())){
      release(&tickslock);
      return -1;
    }
    sleep(&ticks, &tickslock);
  }
  release(&tickslock);
  return 0;
}


#ifdef LAB_PGTBL
int
sys_pgaccess(void)
{
  int page_num;        // 长度
  int mask;            // 掩码
  uint64 buf;
  pagetable_t pagetable_t;

  // 获取当前进程的信息
  pagetable_t = myproc()->pagetable;

  // 获取系统调用参数
  
  // if (argaddr(0, (uint64*)&buf) < 0 || argint(1, (int*)&page_num) < 0 || argaddr(2, (uint64*)&mask) < 0) {
  //   return -1;  // 参数传递错误
  // }
  argaddr(0, &buf);
  argint(1, &page_num);
  argint(2, &mask);
  // 确保不超出上限
  if (page_num > 32) {
    page_num = 32;  // 设置上限
  }

  // 遍历每一页
  for (int i = 0; i < page_num; i++) {
    pte_t *pte = walk(pagetable_t, (buf + i * PGSIZE), 0);

    // 检查访问位 (PTE_A) 是否被设置
    if ((*pte & PTE_A) != 0) {
      *pte = *pte & (~PTE_A);  // 清空 PTE_A 标志位，防止数据干扰
      mask = mask | (1 << i);  // 设置掩码的相应位
    }
  }

  // 将结果掩码复制到用户空间
  if (copyout(pagetable_t, mask, (char*)&mask, sizeof(mask)) < 0) {
    return -1;  // 复制失败
  }

  return 0;  // 系统调用执行成功
}
#endif

uint64
sys_kill(void)
{
  int pid;

  argint(0, &pid);
  return kill(pid);
}

// return how many clock tick interrupts have occurred
// since start.
uint64
sys_uptime(void)
{
  uint xticks;

  acquire(&tickslock);
  xticks = ticks;
  release(&tickslock);
  return xticks;
}

```



```
uint a = m;
uint* b = (uint*)a;
uint* c = a;
```



```
printf("backtrace:\n"); 
  uint64 addr= r_fp();
  uint64* frame =&(((uint64 *)addr)[-1]);
  printf("%d",frame);

  
  uint down = PGROUNDDOWN(addr);
  uint up = PGROUNDUP(addr);
  while(addr > down && addr < up){
    

    printf("%d",*frame);

    addr = frame[-1];

    // &frame表示返回值的地址
    frame = &((uint64 *)frame[-1])[-1];

  }


//  printf("backtrace:\n"); 
//   uint64 fp = r_fp(); // 调用函数来读取当前的帧指针，该函数使用内联汇编来读取s0
  
//   uint64 *frame = (uint64 *)fp; 
//   printf("%d",frame);
//   uint64 up = PGROUNDUP(fp);
//   uint64 down = PGROUNDDOWN(fp); // 计算栈页面的顶部和底部地址
//   while (fp > down && fp < up)
//   {
//     printf("%p\n", frame[-1]);
//     fp = frame[-2];
//     frame = (uint64 *)fp;
//   }

```

-24712*frame:-2147474936down: -28672/nup: -24576/naddr: -24704

-24704up: -24576/ndown: -28672/naddr: -247040x0000000080002208
fp: -24640/n0x000000008000207a
fp: -24608/n0x0000000080001d70

Lab4

```c++
// give up the CPU if this is a timer interrupt.
  if(which_dev == 2){
    if(p->ticks!=0  && !p->is_handler){
      if(p->pre_ticks == p->ticks){
        //保存用户寄存器
        p->trapframe_backup = p->trapframe;
        // 执行中断处理函数
        p->trapframe->epc = (uint64)p->handler;

        p->is_handler = 1;
        p->pre_ticks = 0;
      }
      p->pre_ticks++;
  
    }
    yield();

uint64
sys_sigreturn(void){
  struct proc* p = myproc();
  p->trapframe = p ->trapframe_backup;
  p->is_handler = 0;
  return 0;
}

uint64 
sys_sigalarm(void){
  int ticks;
  uint64 handler;
  argint(0,&ticks);
  argaddr(0,&handler);
  struct proc* p = myproc();
  p->ticks = ticks;
  p->handler = (void*)handler;
  p->pre_ticks = 0;
  return 0;
}

/ Look in the process table for an UNUSED proc.
// If found, initialize state required to run in the kernel,
// and return with p->lock held.
// If there are no free procs, or a memory allocation fails, return 0.
static struct proc*
allocproc(void)
{
  struct proc *p;

  for(p = proc; p < &proc[NPROC]; p++) {
    acquire(&p->lock);
    if(p->state == UNUSED) {
      goto found;
    } else {
      release(&p->lock);
    }
  }
  return 0;

found:
  p->pid = allocpid();
  p->state = USED;

  // Allocate a trapframe page.
  if((p->trapframe = (struct trapframe *)kalloc()) == 0){
    freeproc(p);
    release(&p->lock);
    return 0;
  }

   // Allocate a trapframe backup page.
  if((p->trapframe_backup = (struct trapframe *)kalloc()) == 0){
    freeproc(p);
    release(&p->lock);
    return 0;
  }

  // An empty user page table.
  p->pagetable = proc_pagetable(p);
  if(p->pagetable == 0){
    freeproc(p);
    release(&p->lock);
    return 0;
  }

  // 初始化中断
  p->pre_ticks = 0;
  p->ticks = 0;
  p->handler = 0;
  p->is_handler = 0;

  // Set up new context to start executing at forkret,
  // which returns to user space.
  memset(&p->context, 0, sizeof(p->context));
  p->context.ra = (uint64)forkret;
  p->context.sp = p->kstack + PGSIZE;

  return p;
}

```



```c
// Buffer cache.
//
// The buffer cache is a linked list of buf structures holding
// cached copies of disk block contents.  Caching disk blocks
// in memory reduces the number of disk reads and also provides
// a synchronization point for disk blocks used by multiple processes.
//
// Interface:
// * To get a buffer for a particular disk block, call bread.
// * After changing buffer data, call bwrite to write it to disk.
// * When done with the buffer, call brelse.
// * Do not use the buffer after calling brelse.
// * Only one process at a time can use a buffer,
//     so do not keep them longer than necessary.


#include "types.h"
#include "param.h"
#include "spinlock.h"
#include "sleeplock.h"
#include "riscv.h"
#include "defs.h"
#include "fs.h"
#include "buf.h"

struct {
  struct spinlock lock;
  struct buf buf[BUCKET_NUM];

  // Linked list of all buffers, through prev/next.
  // Sorted by how recently the buffer was used.
  // head.next is most recent, head.prev is least.
  // struct buf head;
} bcache;

void
binit(void)
{
  struct buf *b;

  initlock(&bcache.lock, "bcache");

  for(b = bcache.buf; b < bcache.buf+BUCKET_NUM;b++){
    b->ticks = 0;

  }

  // // Create linked list of buffers
  // bcache.head.prev = &bcache.head;
  // bcache.head.next = &bcache.head;
  // for(b = bcache.buf; b < bcache.buf+NBUF; b++){
  //   b->next = bcache.head.next;
  //   b->prev = &bcache.head;
  //   initsleeplock(&b->lock, "buffer");
  //   bcache.head.next->prev = b;
  //   bcache.head.next = b;
  // }
}

// Look through buffer cache for block on device dev.
// If not found, allocate a buffer.
// In either case, return locked buffer.
static struct buf*
bget(uint dev, uint blockno)
{
  struct buf *b;

  int bucket = blockno % BUCKET_NUM;

  

  // Is the block already cached?
  // for(b = bcache.head.next; b != &bcache.head; b = b->next){
  //   if(b->dev == dev && b->blockno == blockno){
  //     b->refcnt++;
  //     release(&bcache.lock);
  //     acquiresleep(&b->lock);
  //     return b;
  //   }
  // }
  b =  &bcache.buf[bucket];
  for(int i = 0; b; i++){
    acquiresleep(&bcache.buf[bucket+i].lock);
    if(b->dev == dev && b -> blockno == blockno){
      
       b->refcnt++;
       return b;

    }
  }

  // // Not cached.
  // // Recycle the least recently used (LRU) unused buffer.
  // for(b = bcache.head.prev; b != &bcache.head; b = b->prev){
  //   if(b->refcnt == 0) {
  //     b->dev = dev;
  //     b->blockno = blockno;
  //     b->valid = 0;
  //     b->refcnt = 1;
  //     release(&bcache.lock);
  //     acquiresleep(&b->lock);
  //     return b;
  //   }
  // }
  // panic("bget: no buffers");

  uint min_use = __INT32_MAX__;
  struct buf* need_lru = b;

  for(int i = 0; b; i++){
    
    if(&bcache.buf[i].ticks < min_use && b->refcnt == 0){
      acquiresleep(&bcache.buf[i].lock);
      min_use = &bcache.buf[i].ticks;
      need_lru = &bcache.buf[i];
    }
  }
  need_lru->valid =0;
  for(int i = 0; b; i++){
     b = &bcache.buf[bucket+i];
  }
  b->dev = dev;
  b->blockno = blockno;
  b->valid = 0;
  b->refcnt = 1;
  releasesleep(need_lru);
  acquiresleep(b);

  return b;
  

  
}

// Return a locked buf with the contents of the indicated block.
struct buf*
bread(uint dev, uint blockno)
{
  struct buf *b;

  b = bget(dev, blockno);
  if(!b->valid) {
    virtio_disk_rw(b, 0);
    b->valid = 1;
  }
  return b;
}

// Write b's contents to disk.  Must be locked.
void
bwrite(struct buf *b)
{
  if(!holdingsleep(&b->lock))
    panic("bwrite");
  virtio_disk_rw(b, 1);
}

// Release a locked buffer.
// Move to the head of the most-recently-used list.
void
brelse(struct buf *b)
{
  if(!holdingsleep(&b->lock))
    panic("brelse");

  releasesleep(&b->lock);

  acquire(&bcache.lock);
  b->refcnt--;
  if (b->refcnt == 0) {
    // no one is waiting for it.
    b->next->prev = b->prev;
    b->prev->next = b->next;
    b->next = bcache.head.next;
    b->prev = &bcache.head;
    bcache.head.next->prev = b;
    bcache.head.next = b;
  }
  
  release(&bcache.lock);
}

void
bpin(struct buf *b) {
  acquire(&bcache.lock);
  b->refcnt++;
  release(&bcache.lock);
}

void
bunpin(struct buf *b) {
  acquire(&bcache.lock);
  b->refcnt--;
  release(&bcache.lock);
}



```



```
uint64 sys_mmap(void){
    struct proc *p = myproc();
    uint64* addr;
    struct file* fd;
    argaddr(1,&addr);
    int index = (uint64)addr / PGSIZE;
    struct vma vma = p->vmas[index];
    argaddr(2,&vma.len);
    argaddr(3,&vma.flags);
    argaddr(5,&fd);
    argaddr(6,&vma.offset);
    fd->ref++;
    vma.fd = fd;

}
uint64 sys_unmap(void){
   struct proc *p = myproc();
   uint64* addr;
   char* size;
  struct file* fd;
  argaddr(1,&addr);
  argint(1,&size);
  int index = (uint64)addr / PGSIZE;
  struct vma vma = p->vmas[index];
  
  uvmunmap(p->pagetable,addr,size,1);
  vma.addr = addr + (uint64)size;
  vma.len = vma.len - (uint64)size;
  if(vma.len <= 0){
     vma.file->ref--;
  }
   // 写回取消映射页面
   if(vma.flags == MAP_SHARED){
     filewrite(vma.file,vma.addr,size);
   }
}
```

