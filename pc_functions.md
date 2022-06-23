# Methods used for PC

## Pthreads (POSIX threads)
For all these functions, the return value is
- `0` if joining/creating/waiting etc. is successful
- `1` if not
in this case the `errno` is also set to a value specific to what went wrong

### Threads
```c
pthread_create(pthread_t *thread, pthread_attr_t *attr, void* (*start_routine)(void* args), void* args)

```
Arguments
* `thread`            pointer to the thread to create
* `attr`              attributes ((un)detached, (a)synchronous), using NULL will use default attributes (undetached, synchr.)
* `start_routine`     threaded function to call after creating
* `args`              arguments to pass to the threaded function
for a function with multiple arguments, use a struct to collect these
```c
pthread_join(pthread_t *thread, void** status)
```
Arguments
* `thread`            pointer to the thread to join with the calling thread
* `status`            contains a pointer to the return value of the threaded function if there is any

### Mutexes
Used to (un)lock access to section in threaded function
```c
pthread_mutex_init(pthread_mutex_t *mutex, const pthread_mutexattr_t *mutexattr)
pthread_mutex_destroy(pthread_mutex_t *mutex)
```
Dynamically initialize/destroy mutex, otherwise you can use
`pthread_mutex_t mutex = PTHREAD_MUTEX_INITIALIZER`
for default initialization

```c
pthread_mutex_lock(pthread_mutex_t* mutex);
pthread_mutex_unlock(pthread_mutex_t* mutex);
```
Arguments
* `mutex`             mutex to (unlock access to section in threaded function

### Condition variables
Used to synchronize threads by letting them wait on each other
```c
pthread_cond_init(pthread_cond_t *cond, pthread_condattr_t *attr);
pthread_cond_destroy(pthread_cond_t *cond);
```
Dynamically initialize/destroy condition variable, otherwise you can use
`pthread_cond_t cond = PTHREAD_MUTEX_INITIALIZER`
for default initialization

```c
pthread_cond_wait(pthread_cond_t* cond, pthread_mutex_t* mutex)
```
Arguments
* `cond`              condition variable to block on
* `mutex`             mutex that must be locked by calling thread before cond_wait

```c
pthread_cond_signal(pthread_cond_t* cond)
```
Unlocks a single thread that is blocked on the condition variable
If more than one is blocked, the order of unblocking is unspecified
It doesn't do anything if no threads are currently blocked

```c
pthread_cond_broadcast(pthread_cond_t* cond)
```
Unlocks all threads that are blocked on the condition variable
It doesn't do anything if no threads are currently blocked

### Semaphores
Used as a **protected atomic variable** for restricting access to shared resources, e.g. variables. <br>
The initial value is *how many resources* should share this variable

```c
typedef struct {
    pthread_mutex_t mutex;
    pthread_cond_t cond;
    unsigned int value;
} Semaphore;
```
```c
void initSemaphore(Semaphore *s, unsigned int v) {
    pthread_mutex_lock(&(s->mutex));
    s->value = v;
    pthread_mutex_unlock(&(s->mutex));
}
```
```c
void downSemaphore(Semaphore *s) {
    pthread_mutex_lock(&(s->mutex));
    while (s->value == 0) {
        pthread_cond_wait(&(s->cond), &(s->mutex));
    }
    s->value--;
    pthread_mutex_unlock(&(s->mutex));
}
```
```c
void upSemaphore(Semaphore *s) {
    pthread_mutex_lock(&(s->mutex));
    s->value++;
    if (s->value == 1) {
        pthread_cond_broadcast(&(s->cond));
    }
    pthread_mutex_unlock(&(s->mutex));
}
```

### Barriers
Used to synchronize *all* threads by letting them wait on each other at the barrier
```c
pthread_barrier_init(pthread_barrier_t *barrier, const pthread_barrierattr_t *attr, unsigned count); 
pthread_barrier_destroy(pthread_barrier_t *barrier) 
```
Arguments
- `barrier` barrier to initialize/destroy
- `attr` attributes, using NULL will use default
- `count` number of threads that must call `pthread_barrier_wait` before any stop waiting

```c
pthread_barrier_wait(pthread_barrier_t* barrier)
```
- `barrier` barrier variable to determine which threads should wait on each other


## OpenMP (Open Multi-Processing)
Uses `#pragma` compiler directives to handle parallelization, which can be used in C/C++ and Fortran. <br>
If the code is not compiled with an `-fopenmp` or `-openmp` flag, *the directives are ignored and code runs sequentially*.

The number of threads working on a parallel region are determined in this order
1. Evaluation of an `if` clause
```c
#pragma omp parallel if(condition)      // if condition is false, parallalism is prevented
```
2. Setting the `num_threads` clause
```c
#pragma omp parallel num_threads(4)
```
3. With **environment variables**:
```bash
export OMP_NUM_THREADS=2
```
4. Using default number of available cores

### Base methods
```c
#pragma omp parallel
#pragma omp parallel num_threads(...)
#pragma omp parallel shared(...) private(...) firstprivate(...) lastprivate(...)
```
Basis for parallizing any region in code<br>
By default
- the number of threads is used according to the above
- *All* variables in scope are shared among threads<br>
Variables that are created within the parallel section are thread specific however

There are many to change these defaults:
- `num_threads` how many threads should execute the region

As for variable sharing
- `shared(x)` variable `x` is shared among all threads
- `private(x)` each thread gets its own **uninitialized** private variable `x`
- `firstprivate(x)` same as `private` but the initial value of `x` will be that of the original variable
- `lastprivate(x)` same as `private` but the value of the original variable is set to the private variable `x` that is last changed
```c
#pragma omp single
```
Execute region from a single thread instead of all of them
```c
#pragma omp atomic
```
Write to *one variable* from a single thread so race conditions do not occur<br>
Essentially this is a one-line version of the `single` pragma
```c
#pragma omp critical
```
Let threads execute the region one-by-one, i.e. they should wait on each other before entering<br>
This is the equivalent of `pthread_lock(...) code pthread_unlock(...)` for pthreads.
```c
#pragma omp barrier
```
Let all threads wait on each other before continuing
```c
#pragma omp threadprivate(x)
```
Create a thread private variable from a global variable `x`


### For-loops
```c
#pragma omp parallel for
#pragma omp parallel for    nowait ordered collapse(...)
```
Basic way to parallelise a for-loop<br>
The number of iterations will be equally divided among the number of threads<br>
Some control options
- `nowait` There is an implicit barrier after the body of the for-loop<br>
Setting `nowait` removes this barrier, i.e. threads go on after they finish their part of the loop
- `ordered` Specify this if the iterations of the loop should be executed in order, e.g.
```c
for (int i = 1; i < n; i++)
    a[i] = a[i-1] + 1;
```
needs to be in order because `i` depends on `i-1`
- `collapse(n)` Simplifies a `n` layer for-loop into a single one<br>
This is useful when e.g. there are way more threads than iterations for the outer loop<br>
**This only work if the loops are independent of each other**
```c
#pragma omp for collapse(2)
for (int i = 0; i < 4; i++) {           // without collapse, only 4 threads are used for parallelisation
    for (int j = 0; j < 100; j++) {
        ...
    }
}

// collapsing yields this
for (int i = 0; i < 400; i++) {         // now as many as 400 threads can be used for parallelisation
    ...
}

```



### Sections
```c
#pragma omp sections
#pragma omp sections nowait
```
Specifies a parallel region in which different threads all execute a `section` <br>
There is an *implicit barrier* at the end of this region, i.e. all sections should be completed before any thread continues <br>
Use `nowait` if this barrier should not exist
```c
#pragma omp section
```
Section that should be executed by one thread

#### Example
```c
#pragma omp sections
{
    #pragma omp section
    // do thing 1

    #pragma omp section
    // do thing 2
}
```

### Tasks
```c
#pragma omp task
```
Enqueues (**but not executes!**) a task to the task queue<br>
*This is why this pragma should only be used inside of a `#pragma omp single` region*<br>
```c
#pragma omp task        // wrong

#pragma omp single
{
    #pragma omp task    // right
}
```
Whenever a thread is available, it will pop the task from the queue and execute the code section 
```c
#pragma omp taskwait
```
There is *no implicit barrier* for tasks, so threads can continue before all tasks are completed<br>
If you want all tasks to be completed before the next section is executed, use this pragma.

### Reduce
This can be applied to a piece of code to collect results from all threads.
The most common is for `reduce` to be used with a `for loop`:
```c
int s = 0
#pragma omp parallel for reduction(+:s)
for (int i = 0; j < n; j++)
    s += a[j];
```
This achieves the following
- Each threads creates its own `s` to computes its part of the global sum
- All thread `s` variables are reduced to a global variable `s` (which should be declared before!)


## MPI (Message Passing Interface)
In constrast to OpenMP, MPI *doesn't have a global address space / shared variables*<br>
This implies each process has *its own local address space* that it works with.
To provide communication between process, message are *sent* and *received*.

For all the MPI functions, the return value is
- `0`       if function is successful
- `not 0`   if not, then a specific error code is returned

### Setup
```c
MPI_Init(&argc, *argv)
// MPI parallelised code
MPI_Finalize()
```
Begin/start point of any parallelisation with MPI<br>
*All* MPI calls need to be inbetween these two


### Communicators and ranks
In MPI, *communicator*s are groups of processes that can communicate with each other<br>
The default communicator `MPI_COMM_WORLD` consists of all processes<br>
```c
MPI_Comm_size(MPI_Comm comm, int* size)
```
then tells the current process what the size of of the `comm` group is by storing it in `size`.
```c
MPI_Comm_rank(MPI_Comm comm, int* rank)
```
In the same way, a process can verify its *rank*, which is its id within the communicator it belongs to by calling `MPI_Comm_rank`.

#### Example (with global communicator)
```c
int groupSize;      // how many processors share the global communicator
int rank;           // which id this process has within the global communicator

// Sets the groupSize / rank for the current process to the right value
MPI_Comm_size(MPI_COMM_WORLD, &groupSize);
MPI_Comm_rank(MPI_COMM_WORLD, &rank);
```

### Group and custom communicators
Say you want to create your own `MPI_Comm` communicator instead of using the global one<br>
You can do so by taking a *group* (think subset) from another communicator to do so. <br>
The process would go like this:
0. Create variables for new groups and communicator
```c
MPI_Group orig_group;       // every process in MPI_COMM_WORLD
MPI_Group new_group;        // every process in the group for your communicator
MPI_Comm new_comm;          // communicator you are going to create
```
1. Collect the group of all processors within a communicator
```c
MPI_Comm_group(MPI_COMM_WORLD, &orig_group);
```
2. Make subgroup that will represent all processors within your new communicator
```c
int nProcs;                 // how many processors there are in the original group
int ranks1[], rank2[]       // separate processors into two groups
if (rank < nProcs)
    MPI_Group_incl(orig_group, nProcs, ranks1, &new_group);     // put processor 0, 1 etc. in first group
else
    MPI_Group_incl(orig_group, nProcs, ranks2, &new_group);     // put processor nProcs/2, ... into second group
```
3. Create the new communicator based on the subgroup you want
```c
MPI_Comm_create(MPI_COMM_WORLD, new_group, &new_comm)
```

As a result, you now have two communicators which both include half of the processors in `MPI_COMM_WORLD`

### Passing around messages
#### Sending (to one)
The way processes communicate with each other is to send/receive messages<br>
There are a number of different ways of doing so, e.g. for sending you have
- `Send` Blocking method that waits until the process has sent the message
- `Ssend` Same as `Send`, but it *synchronizes* with receiving<br>
So when this function returns, the destination has started receiving the method
- `ISend` Non-blocking version of `Send`, returns immediately<br>
You have to now verify yourself that the `msg` buffer can be reused, because the actual send might not have completed
- `BSend` Same as `Send`, but it returns as soon as the buffer is safe to be reused

The prototype for all these functions is the same, which is
```c
MPI_Send(void* msg, int cnt, MPI_<Datatype> type, int dest, int tag, MPI_Comm comm)
```
Arguments
- `msg` Starting address of the data to send
- `cnt` How many data items to send (e.g. send the first 3 letters of a word in a `char[]`)
- `type` What data type is sent <br>
OpenMPI provides most basic types, such as `MPI_INT`, `MPI_CHAR` etc.
- `dest` Rank of the process within the `comm` group that the message is sent to
- `tag` Message tag to specify details about the sending process (e.g. rank of sending process)
- `comm` Communicator that is addressed in the message

So lets say you want to sent the 3rd letter from a `char* word` to the 4th process in the global group<br>
This can be achieved with
```c
MPI_Send(msg + 2, 1, MPI_CHAR, 3, 99, MPI_COMM_WORLD)
```

#### Receiving
For receiving, the prototype is very similar to that of sending:
```c
MPI_Recv(void* msg, int cnt, MPI_<Datatype> type, int src, int tag, MPI_Comm comm, MPI_Status* status)
```
Note that
- Instead of specifying the destination `dest`, we now specify the source process `src` from which the message is expected to be sent
- The additional paramater `status` stores the status of receiving, which includes e.g. how many items where received and whether the receiving was cancelled or not

Just like with sending, `IRecv` exists which doesn't wait for the message to be fully received.
We then need to use the `Wait` method to wait until this has been fully completed:
```c
MPI_Wait(MPI_Request* request, MPI_Status* status)
```
This function waits until the `request` is completed and saves the request information in `status`.

```c
MPI_Test(MPI_Request* request, int* flag, MPI_Status* status)
```
Additionally, there is also a method for testing whether a message has been received fully or not<br>
The `flag` is then set to a non-zero value to indicate the request is completed


#### Broadcasting
For sending a message from one process to all processes in the communicator, use the `Broadcast` function:
```c
MPI_Bcast(void* msg, int count, MPI_<DataType> type, int root, MPI_Comm comm)
```
Arguments
- `msg` Starting address of the data to send for `root` process
- `cnt` How many data items to send
- `type` What data type is sent
- `root` Rank of the sending process within the communicator
- `comm` Communicator that is addressed in the message

**Note that all processes need to call this function to receive the message, not just the one sending it!**

### Distributing data
Above we mentioned ways to send data from one process to another<br>
There are also different ways of sending/receiving data between *all* processors within a communicator<br>
These are

| Operation | Descr. | Processes Data (before) | Processes Data (after) |
| :--- |:---| :-- | ---:|
| `Scatter` | Divide data equally over processes | P0: [A B C]<br>P1:[]<br>P2: [] | P0: [A]<br> P1: [B]<br>P2: [C] |
| `Gather` | Reverse of `Scatter`, gather data in first process | P0: [A]<br>P1:[B]<br>P2: [C] | P0: [A B C]<br> P1: []<br>P2: [] |
| `AllGather` | `Gather` in each process | P0: [A]<br>P1:[B]<br>P2: [C] | P0: [A B C]<br> P1: [A B C]<br>P2: [A B C] |
| `Alltoall` | "Transpose" data, or in other words<br>Collect `i`th elem from each process in process `i` | P0: [A D G]<br>P1:[B E H]<br>P2: [C F I] | P0: [A B C]<br> P1: [D E F]<br>P2: [G H I] |

The function prototype for all these methods is the same:
```c
MPI_Scatter(void* src, int srcCnt, MPI_<Datatype> srcType, void* dest, int destCnt, MPI_<Datatype> destType, int root, MPI_Comm comm)
```
Arguments
- `src` Starting address of the data to send
- `dest` Starting address of where data should be received
- `srcCnt / destCnt` How many data items should be sent/received
- `srcType / destType` What data type is sent/as what type the data should be received
- `root` Rank of the sending process within the communicator
- `comm` Communicator that is addressed in the message

### Reduce/Scan
Finally, MPI is also able to do reduce/scan results between processors<br>
Their prototypes are really similar:
```c
MPI_Reduce(void* src, void* dest, int count, MPI_<DataType> type, MPI_<Op> op, int root, MPI_Comm comm);
MPI_Scan(void* src, void* dest, int count, MPI_<DataType> type, MPI_<Op> op, MPI_Comm comm);
```
Arguments
- `src` Starting address of the data to reduce/scan
- `dest` Starting address of where reduced/scanned data should be stored<br>
**Note that `src` and `dest` should not overlap, this is gonna give errors!**
- `type` What data type is reduced/scanned
- `op` How the data is reduced/scanned<br>
Just like with `MPI_<DataType>`, MPI includes pretty much every operator type you'd expect<br>
Examples are `MPI_MAX`, `MPI_SUM`, `MPI_PROD` and so on
- `comm` Communicator that specifies all processors to reduce/scan over

The difference is that `Reduce` also needs to specify a `root` process in which the result will be stored.