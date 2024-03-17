#include <linux/module.h> // THIS_MODULE, MODULE_VERSION, ...
#include <linux/init.h>   // module_{init,exit}
#include <linux/proc_fs.h>
#include <linux/sched/signal.h> // for_each_process()
#include <linux/seq_file.h>
#include <linux/fs.h>
#include <linux/sched.h>
#include <linux/mm.h> // get_mm_rss()
#include <linux/kernel.h>
#include <asm/uaccess.h>

MODULE_LICENSE("GPL");
MODULE_AUTHOR("DannyT8355");
MODULE_DESCRIPTION("Informacion cpu");
MODULE_VERSION("1.0");

struct task_struct *task;       // sched.h para tareas/procesos
struct task_struct *task_child; // index de tareas secundarias
struct list_head *list;         // lista de cada tareas
struct sysinfo si;

static void init_meminfo(void) {
    si_meminfo(&si);
}

static int escribir_a_proc(struct seq_file *file_proc, void *v)
{
    int running = 0;
    int sleeping = 0;
    int zombie = 0;
    int stopped = 0;
    unsigned long rss;
    unsigned long total_ram_pages;
    
    
    total_ram_pages = totalram_pages();
    if (!total_ram_pages) {
        pr_err("No memory available\n");
        return -EINVAL;
    }
    
    #ifndef CONFIG_MMU
        pr_err("No MMU, cannot calculate RSS.\n");
        return -EINVAL;
    #endif
    
    unsigned long total_cpu_time = jiffies_to_msecs(get_jiffies_64());
    unsigned long total_usage = 0;

    for_each_process(task) {
        unsigned long cpu_time = jiffies_to_msecs(task->utime + task->stime);
        total_usage += cpu_time;
    }
    init_meminfo();
    //---------------------------------------------------------------------------
    seq_printf(file_proc, "{\n\t\"totalram\": %lu,", si.totalram);
    seq_printf(file_proc, "\n\t\"freeram\": %lu,", si.freeram);
    seq_printf(file_proc, "\n\t\"cpu_total\":%ld,", total_cpu_time);
    seq_printf(file_proc, "\n\t\"cpu_porcentaje\":%ld,", (total_usage * 100) / total_cpu_time);
    seq_printf(file_proc, "\n\t\"processes\": [");
    int b = 0;

    for_each_process(task)
    {
        if (task->mm)
        {
            rss = get_mm_rss(task->mm) << PAGE_SHIFT;
        }
        else
        {
            rss = 0;
        }
        if (b == 0)
        {
            seq_printf(file_proc, "\n\t\t{");
            b = 1;
        }
        else
        {
            seq_printf(file_proc, ",\n\t\t{");
        }
        seq_printf(file_proc, "\n\t\t\t\"pid\":%d,", task->pid);
        seq_printf(file_proc, "\n\t\t\t\"name\":\"%s\",", task->comm);
        seq_printf(file_proc, "\n\t\t\t\"user\": %u,", task->cred->uid);
        seq_printf(file_proc, "\n\t\t\t\"state\":%d,", task->__state);
        int porcentaje = (rss * 100) / total_ram_pages;
        seq_printf(file_proc, "\n\t\t\t\"ram\":%d,", porcentaje);

        seq_printf(file_proc, "\n\t\t\t\"child\": [");
        int a = 0;
        list_for_each(list, &(task->children))
        {
            task_child = list_entry(list, struct task_struct, sibling);
            if (a != 0)
            {
                seq_printf(file_proc, ",\n\t\t\t\t{");
                seq_printf(file_proc, "\n\t\t\t\t\t\"pid\":%d,", task_child->pid);
                seq_printf(file_proc, "\n\t\t\t\t\t\"name\":\"%s\",", task_child->comm);
                seq_printf(file_proc, "\n\t\t\t\t\t\"state\":%d,", task_child->__state);
                seq_printf(file_proc, "\n\t\t\t\t\t\"pidPadre\":%d", task->pid);
                seq_printf(file_proc, "\n\t\t\t\t}");
            }
            else
            {
                seq_printf(file_proc, "\n\t\t\t\t{");
                seq_printf(file_proc, "\n\t\t\t\t\t\"pid\":%d,", task_child->pid);
                seq_printf(file_proc, "\n\t\t\t\t\t\"name\":\"%s\",", task_child->comm);
                seq_printf(file_proc, "\n\t\t\t\t\t\"state\":%d,", task_child->__state);
                seq_printf(file_proc, "\n\t\t\t\t\t\"pidPadre\":%d", task->pid);
                seq_printf(file_proc, "\n\t\t\t\t}");
                a = 1;
            }
        }
        if (a != 0) {
            seq_printf(file_proc, "\n\t\t\t]");
        } else {
            seq_printf(file_proc, "]");
        }
        a = 0;

        if (task->__state == 0)
        {
            running += 1;
        }
        else if (task->__state == 1)
        {
            sleeping += 1;
        }
        else if (task->__state == 4)
        {
            zombie += 1;
        }
        else
        {
            stopped += 1;
        }
        seq_printf(file_proc, "\n\t\t}");
    }
    b = 0;
    seq_printf(file_proc, "\n\t],");
    seq_printf(file_proc, "\n\t\"running\":%d,", running);
    seq_printf(file_proc, "\n\t\"sleeping\":%d,", sleeping);
    seq_printf(file_proc, "\n\t\"zombie\":%d,", zombie);
    seq_printf(file_proc, "\n\t\"stopped\":%d,", stopped);
    seq_printf(file_proc, "\n\t\"total\":%d", running + sleeping + zombie + stopped);
    seq_printf(file_proc, "\n}");
    return 0;
}

static int abrir_aproc(struct inode *inode, struct file *file)
{
    return single_open(file, escribir_a_proc, NULL);
}

static struct proc_ops archivo_operaciones = {
    .proc_open = abrir_aproc,
    .proc_read = seq_read
};

static int __init modulo_init(void)
{
    proc_create("ram_cpu", 0, NULL, &archivo_operaciones);
    printk(KERN_INFO "Insertar Modulo CPU\n");
    return 0;
}

static void __exit modulo_cleanup(void)
{
    remove_proc_entry("ram_cpu", NULL);
    printk(KERN_INFO "Remover Modulo CPU\n");
}

module_init(modulo_init);
module_exit(modulo_cleanup);