	.file	"fe-connect.c"
	.text
	.p2align 4,,15
	.type	pwdfMatchesString, @function
pwdfMatchesString:
.LFB1002:
	.cfi_startproc
	testq	%rdi, %rdi
	je	.L11
	testq	%rsi, %rsi
	je	.L11
	movzbl	(%rdi), %eax
	cmpb	$42, %al
	je	.L20
	testb	%al, %al
	jne	.L8
	jmp	.L11
	.p2align 4,,10
	.p2align 3
.L5:
	movzbl	1(%rdi), %ecx
	leaq	1(%rdi), %rax
.L6:
	testb	%dl, %dl
	je	.L11
	cmpb	%dl, %cl
	jne	.L11
	leaq	1(%rax), %rdi
	movzbl	1(%rax), %eax
	addq	$1, %rsi
	testb	%al, %al
	je	.L11
.L8:
	cmpb	$92, %al
	movzbl	(%rsi), %edx
	je	.L5
	movzbl	(%rdi), %ecx
	movq	%rdi, %rax
	cmpb	$58, %cl
	jne	.L6
	testb	%dl, %dl
	jne	.L6
	leaq	1(%rdi), %rax
	ret
	.p2align 4,,10
	.p2align 3
.L11:
	xorl	%eax, %eax
	ret
	.p2align 4,,10
	.p2align 3
.L20:
	cmpb	$58, 1(%rdi)
	jne	.L8
	leaq	2(%rdi), %rax
	ret
	.cfi_endproc
.LFE1002:
	.size	pwdfMatchesString, .-pwdfMatchesString
	.section	.rodata.str1.1,"aMS",@progbits,1
.LC0:
	.string	"%s"
	.text
	.p2align 4,,15
	.type	defaultNoticeProcessor, @function
defaultNoticeProcessor:
.LFB1001:
	.cfi_startproc
	movq	stderr@GOTPCREL(%rip), %rax
	movq	%rsi, %rdx
	leaq	.LC0(%rip), %rsi
	movq	(%rax), %rdi
	xorl	%eax, %eax
	jmp	pg_fprintf@PLT
	.cfi_endproc
.LFE1001:
	.size	defaultNoticeProcessor, .-defaultNoticeProcessor
	.section	.rodata.str1.8,"aMS",@progbits,1
	.align 8
.LC1:
	.string	"invalid integer value \"%s\" for connection option \"%s\"\n"
	.text
	.p2align 4,,15
	.type	parse_int_param, @function
parse_int_param:
.LFB928:
	.cfi_startproc
	pushq	%r15
	.cfi_def_cfa_offset 16
	.cfi_offset 15, -16
	pushq	%r14
	.cfi_def_cfa_offset 24
	.cfi_offset 14, -24
	movq	%rcx, %r15
	pushq	%r13
	.cfi_def_cfa_offset 32
	.cfi_offset 13, -32
	pushq	%r12
	.cfi_def_cfa_offset 40
	.cfi_offset 12, -40
	movq	%rsi, %r14
	pushq	%rbp
	.cfi_def_cfa_offset 48
	.cfi_offset 6, -48
	pushq	%rbx
	.cfi_def_cfa_offset 56
	.cfi_offset 3, -56
	movq	%rdi, %r12
	movq	%rdx, %r13
	subq	$40, %rsp
	.cfi_def_cfa_offset 96
	movl	$0, (%rsi)
	movq	%fs:40, %rax
	movq	%rax, 24(%rsp)
	xorl	%eax, %eax
	call	__errno_location@PLT
	leaq	16(%rsp), %rsi
	movl	$0, (%rax)
	movl	$10, %edx
	movq	%r12, %rdi
	movq	%rax, %rbp
	call	strtol@PLT
	movq	16(%rsp), %rbx
	cmpq	%r12, %rbx
	je	.L25
	movl	0(%rbp), %edx
	testl	%edx, %edx
	je	.L37
.L25:
	leaq	928(%r13), %rdi
	leaq	.LC1(%rip), %rsi
	xorl	%eax, %eax
	movq	%r15, %rcx
	movq	%r12, %rdx
	call	appendPQExpBuffer@PLT
	xorl	%eax, %eax
.L22:
	movq	24(%rsp), %rbx
	xorq	%fs:40, %rbx
	jne	.L38
	addq	$40, %rsp
	.cfi_remember_state
	.cfi_def_cfa_offset 56
	popq	%rbx
	.cfi_def_cfa_offset 48
	popq	%rbp
	.cfi_def_cfa_offset 40
	popq	%r12
	.cfi_def_cfa_offset 32
	popq	%r13
	.cfi_def_cfa_offset 24
	popq	%r14
	.cfi_def_cfa_offset 16
	popq	%r15
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L37:
	.cfi_restore_state
	movq	%rax, %rdx
	cltq
	cmpq	%rdx, %rax
	jne	.L25
	movzbl	(%rbx), %ebp
	testb	%bpl, %bpl
	je	.L26
	movq	%rdx, 8(%rsp)
	call	__ctype_b_loc@PLT
	movq	8(%rsp), %rdx
	movq	(%rax), %rcx
	leaq	1(%rbx), %rax
	jmp	.L27
	.p2align 4,,10
	.p2align 3
.L28:
	movq	%rax, 16(%rsp)
	addq	$1, %rax
	movzbl	-1(%rax), %ebp
	testb	%bpl, %bpl
	je	.L26
.L27:
	testb	$32, 1(%rcx,%rbp,2)
	jne	.L28
	jmp	.L25
	.p2align 4,,10
	.p2align 3
.L26:
	movl	%edx, (%r14)
	movl	$1, %eax
	jmp	.L22
.L38:
	call	__stack_chk_fail@PLT
	.cfi_endproc
.LFE928:
	.size	parse_int_param, .-parse_int_param
	.p2align 4,,15
	.type	conninfo_getval, @function
conninfo_getval:
.LFB969:
	.cfi_startproc
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset 6, -16
	pushq	%rbx
	.cfi_def_cfa_offset 24
	.cfi_offset 3, -24
	movq	%rdi, %rbx
	subq	$8, %rsp
	.cfi_def_cfa_offset 32
	movq	(%rdi), %rdi
	testq	%rdi, %rdi
	je	.L43
	movq	%rsi, %rbp
	jmp	.L42
	.p2align 4,,10
	.p2align 3
.L46:
	addq	$56, %rbx
	movq	(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L43
.L42:
	movq	%rbp, %rsi
	call	strcmp@PLT
	testl	%eax, %eax
	jne	.L46
	movq	24(%rbx), %rax
	addq	$8, %rsp
	.cfi_remember_state
	.cfi_def_cfa_offset 24
	popq	%rbx
	.cfi_def_cfa_offset 16
	popq	%rbp
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L43:
	.cfi_restore_state
	addq	$8, %rsp
	.cfi_def_cfa_offset 24
	xorl	%eax, %eax
	popq	%rbx
	.cfi_def_cfa_offset 16
	popq	%rbp
	.cfi_def_cfa_offset 8
	ret
	.cfi_endproc
.LFE969:
	.size	conninfo_getval, .-conninfo_getval
	.p2align 4,,15
	.type	defaultNoticeReceiver, @function
defaultNoticeReceiver:
.LFB1000:
	.cfi_startproc
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset 6, -16
	pushq	%rbx
	.cfi_def_cfa_offset 24
	.cfi_offset 3, -24
	subq	$8, %rsp
	.cfi_def_cfa_offset 32
	movq	128(%rsi), %rbp
	testq	%rbp, %rbp
	je	.L47
	movq	%rsi, %rbx
	movq	%rsi, %rdi
	call	PQresultErrorMessage@PLT
	movq	136(%rbx), %rdi
	addq	$8, %rsp
	.cfi_remember_state
	.cfi_def_cfa_offset 24
	movq	%rax, %rsi
	popq	%rbx
	.cfi_def_cfa_offset 16
	movq	%rbp, %rax
	popq	%rbp
	.cfi_def_cfa_offset 8
	jmp	*%rax
	.p2align 4,,10
	.p2align 3
.L47:
	.cfi_restore_state
	addq	$8, %rsp
	.cfi_def_cfa_offset 24
	popq	%rbx
	.cfi_def_cfa_offset 16
	popq	%rbp
	.cfi_def_cfa_offset 8
	ret
	.cfi_endproc
.LFE1000:
	.size	defaultNoticeReceiver, .-defaultNoticeReceiver
	.p2align 4,,15
	.type	pqDropServerData, @function
pqDropServerData:
.LFB910:
	.cfi_startproc
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset 6, -16
	pushq	%rbx
	.cfi_def_cfa_offset 24
	.cfi_offset 3, -24
	movq	%rdi, %rbx
	subq	$8, %rsp
	.cfi_def_cfa_offset 32
	movq	376(%rdi), %rdi
	testq	%rdi, %rdi
	je	.L51
	.p2align 4,,10
	.p2align 3
.L52:
	movq	24(%rdi), %rbp
	call	free@PLT
	testq	%rbp, %rbp
	movq	%rbp, %rdi
	jne	.L52
.L51:
	movq	776(%rbx), %rdi
	movq	$0, 384(%rbx)
	movq	$0, 376(%rbx)
	testq	%rdi, %rdi
	je	.L53
	.p2align 4,,10
	.p2align 3
.L54:
	movq	(%rdi), %rbp
	call	free@PLT
	testq	%rbp, %rbp
	movq	%rbp, %rdi
	jne	.L54
.L53:
	movq	800(%rbx), %rdi
	movq	$0, 776(%rbx)
	movl	$0, 784(%rbx)
	movb	$0, 788(%rbx)
	movl	$0, 700(%rbx)
	testq	%rdi, %rdi
	je	.L55
	call	free@PLT
.L55:
	movq	712(%rbx), %rdi
	movq	$0, 800(%rbx)
	movb	$0, 360(%rbx)
	movb	$0, 704(%rbx)
	movb	$0, 705(%rbx)
	movb	$0, 708(%rbx)
	testq	%rdi, %rdi
	je	.L56
	call	free@PLT
.L56:
	movq	$0, 712(%rbx)
	movq	$0, 764(%rbx)
	addq	$8, %rsp
	.cfi_def_cfa_offset 24
	popq	%rbx
	.cfi_def_cfa_offset 16
	popq	%rbp
	.cfi_def_cfa_offset 8
	ret
	.cfi_endproc
.LFE910:
	.size	pqDropServerData, .-pqDropServerData
	.p2align 4,,15
	.type	freePGconn, @function
freePGconn:
.LFB940:
	.cfi_startproc
	pushq	%r14
	.cfi_def_cfa_offset 16
	.cfi_offset 14, -16
	pushq	%r13
	.cfi_def_cfa_offset 24
	.cfi_offset 13, -24
	pushq	%r12
	.cfi_def_cfa_offset 32
	.cfi_offset 12, -32
	pushq	%rbp
	.cfi_def_cfa_offset 40
	.cfi_offset 6, -40
	pushq	%rbx
	.cfi_def_cfa_offset 48
	.cfi_offset 3, -48
	movq	%rdi, %rbx
	subq	$16, %rsp
	.cfi_def_cfa_offset 64
	movl	328(%rdi), %edx
	movq	%fs:40, %rax
	movq	%rax, 8(%rsp)
	xorl	%eax, %eax
	testl	%edx, %edx
	jle	.L73
	xorl	%r12d, %r12d
	movq	%rsp, %r13
	.p2align 4,,10
	.p2align 3
.L74:
	movslq	%r12d, %rax
	movq	%r13, %rsi
	movl	$2, %edi
	leaq	(%rax,%rax,4), %rbp
	movq	320(%rbx), %rax
	movq	%rbx, (%rsp)
	addl	$1, %r12d
	salq	$3, %rbp
	addq	%rbp, %rax
	movq	16(%rax), %rdx
	call	*(%rax)
	movq	320(%rbx), %rax
	movq	8(%rax,%rbp), %rdi
	call	free@PLT
	cmpl	%r12d, 328(%rbx)
	jg	.L74
.L73:
	movq	400(%rbx), %rdx
	testq	%rdx, %rdx
	je	.L75
	movl	392(%rbx), %eax
	testl	%eax, %eax
	jle	.L76
	xorl	%r12d, %r12d
	movq	$-1, %r13
	.p2align 4,,10
	.p2align 3
.L81:
	movslq	%r12d, %rax
	leaq	(%rax,%rax,4), %rbp
	salq	$3, %rbp
	leaq	(%rdx,%rbp), %rax
	movq	8(%rax), %rdi
	testq	%rdi, %rdi
	je	.L77
	call	free@PLT
	movq	400(%rbx), %rdx
	leaq	(%rdx,%rbp), %rax
.L77:
	movq	16(%rax), %rdi
	testq	%rdi, %rdi
	je	.L78
	call	free@PLT
	movq	400(%rbx), %rdx
	leaq	(%rdx,%rbp), %rax
.L78:
	movq	24(%rax), %rdi
	testq	%rdi, %rdi
	je	.L79
	call	free@PLT
	movq	400(%rbx), %rdx
	leaq	(%rdx,%rbp), %rax
.L79:
	movq	32(%rax), %r14
	testq	%r14, %r14
	je	.L80
	movq	%r14, %rdi
	call	strlen@PLT
	movq	%r13, %rdx
	movq	%rax, %rsi
	movq	%r14, %rdi
	call	__explicit_bzero_chk@PLT
	movq	400(%rbx), %rax
	movq	32(%rax,%rbp), %rdi
	call	free@PLT
	movq	400(%rbx), %rdx
.L80:
	addl	$1, %r12d
	cmpl	%r12d, 392(%rbx)
	jg	.L81
.L76:
	movq	%rdx, %rdi
	call	free@PLT
.L75:
	movq	(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L82
	call	free@PLT
.L82:
	movq	56(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L83
	call	free@PLT
.L83:
	movq	320(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L84
	call	free@PLT
.L84:
	movq	8(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L85
	call	free@PLT
.L85:
	movq	16(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L86
	call	free@PLT
.L86:
	movq	24(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L87
	call	free@PLT
.L87:
	movq	32(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L88
	call	free@PLT
.L88:
	movq	40(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L89
	call	free@PLT
.L89:
	movq	48(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L90
	call	free@PLT
.L90:
	movq	64(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L91
	call	free@PLT
.L91:
	movq	72(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L92
	call	free@PLT
.L92:
	movq	80(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L93
	call	free@PLT
.L93:
	movq	88(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L94
	call	free@PLT
.L94:
	movq	96(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L95
	call	free@PLT
.L95:
	movq	104(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L96
	call	free@PLT
.L96:
	movq	112(%rbx), %rbp
	testq	%rbp, %rbp
	je	.L97
	movq	%rbp, %rdi
	call	strlen@PLT
	movq	%rbp, %rdi
	movq	$-1, %rdx
	movq	%rax, %rsi
	call	__explicit_bzero_chk@PLT
	movq	112(%rbx), %rdi
	call	free@PLT
.L97:
	movq	120(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L98
	call	free@PLT
.L98:
	movq	128(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L99
	call	free@PLT
.L99:
	movq	136(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L100
	call	free@PLT
.L100:
	movq	144(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L101
	call	free@PLT
.L101:
	movq	152(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L102
	call	free@PLT
.L102:
	movq	160(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L103
	call	free@PLT
.L103:
	movq	168(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L104
	call	free@PLT
.L104:
	movq	192(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L105
	call	free@PLT
.L105:
	movq	184(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L106
	call	free@PLT
.L106:
	movq	200(%rbx), %rbp
	testq	%rbp, %rbp
	je	.L107
	movq	%rbp, %rdi
	call	strlen@PLT
	movq	%rbp, %rdi
	movq	$-1, %rdx
	movq	%rax, %rsi
	call	__explicit_bzero_chk@PLT
	movq	200(%rbx), %rdi
	call	free@PLT
.L107:
	movq	208(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L108
	call	free@PLT
.L108:
	movq	216(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L109
	call	free@PLT
.L109:
	movq	176(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L110
	call	free@PLT
.L110:
	movq	224(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L111
	call	free@PLT
.L111:
	movq	256(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L112
	call	free@PLT
.L112:
	movq	264(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L113
	call	free@PLT
.L113:
	movq	232(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L114
	call	free@PLT
.L114:
	movq	240(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L115
	call	free@PLT
.L115:
	movq	248(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L116
	call	free@PLT
.L116:
	movq	408(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L117
	call	free@PLT
.L117:
	movq	352(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L118
	call	free@PLT
.L118:
	movq	712(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L119
	call	free@PLT
.L119:
	movq	808(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L120
	call	free@PLT
.L120:
	movq	832(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L121
	call	free@PLT
.L121:
	movq	856(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L122
	call	free@PLT
.L122:
	movq	272(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L123
	call	free@PLT
.L123:
	leaq	928(%rbx), %rdi
	call	termPQExpBuffer@PLT
	leaq	952(%rbx), %rdi
	call	termPQExpBuffer@PLT
	movq	8(%rsp), %rax
	xorq	%fs:40, %rax
	jne	.L270
	addq	$16, %rsp
	.cfi_remember_state
	.cfi_def_cfa_offset 48
	movq	%rbx, %rdi
	popq	%rbx
	.cfi_def_cfa_offset 40
	popq	%rbp
	.cfi_def_cfa_offset 32
	popq	%r12
	.cfi_def_cfa_offset 24
	popq	%r13
	.cfi_def_cfa_offset 16
	popq	%r14
	.cfi_def_cfa_offset 8
	jmp	free@PLT
.L270:
	.cfi_restore_state
	call	__stack_chk_fail@PLT
	.cfi_endproc
.LFE940:
	.size	freePGconn, .-freePGconn
	.p2align 4,,15
	.type	makeEmptyPGconn, @function
makeEmptyPGconn:
.LFB939:
	.cfi_startproc
	pushq	%rbx
	.cfi_def_cfa_offset 16
	.cfi_offset 3, -16
	movl	$976, %edi
	call	malloc@PLT
	testq	%rax, %rax
	movq	%rax, %rbx
	je	.L271
	leaq	976(%rax), %rdx
	.p2align 4,,10
	.p2align 3
.L273:
	addq	$8, %rax
	movq	$0, -8(%rax)
	cmpq	%rax, %rdx
	jne	.L273
	leaq	defaultNoticeReceiver(%rip), %rax
	movl	$16384, %edi
	movl	$1, 336(%rbx)
	movl	$0, 340(%rbx)
	movl	$0, 344(%rbx)
	movq	%rax, 288(%rbx)
	leaq	defaultNoticeProcessor(%rip), %rax
	movb	$0, 366(%rbx)
	movb	$0, 367(%rbx)
	movl	$8, 748(%rbx)
	movq	%rax, 304(%rbx)
	movl	$0, 784(%rbx)
	movb	$0, 788(%rbx)
	movl	$1, 792(%rbx)
	movl	$1, 796(%rbx)
	movl	$-1, 416(%rbx)
	movl	$16384, 816(%rbx)
	call	malloc@PLT
	movl	$16384, %edi
	movq	%rax, 808(%rbx)
	movl	$16384, 840(%rbx)
	call	malloc@PLT
	movl	$512, %edi
	movq	%rax, 832(%rbx)
	movl	$32, 864(%rbx)
	call	malloc@PLT
	leaq	928(%rbx), %rdi
	movq	%rax, 856(%rbx)
	call	initPQExpBuffer@PLT
	leaq	952(%rbx), %rdi
	call	initPQExpBuffer@PLT
	cmpq	$0, 808(%rbx)
	je	.L274
	cmpq	$0, 832(%rbx)
	je	.L274
	cmpq	$0, 856(%rbx)
	je	.L274
	cmpq	$0, 944(%rbx)
	je	.L274
	cmpq	$0, 968(%rbx)
	je	.L274
.L271:
	movq	%rbx, %rax
	popq	%rbx
	.cfi_remember_state
	.cfi_def_cfa_offset 8
	ret
.L274:
	.cfi_restore_state
	movq	%rbx, %rdi
	xorl	%ebx, %ebx
	call	freePGconn
	movq	%rbx, %rax
	popq	%rbx
	.cfi_def_cfa_offset 8
	ret
	.cfi_endproc
.LFE939:
	.size	makeEmptyPGconn, .-makeEmptyPGconn
	.section	.rodata.str1.1
.LC2:
	.string	"service"
.LC3:
	.string	"out of memory\n"
	.text
	.p2align 4,,15
	.type	fillPGconn, @function
fillPGconn:
.LFB917:
	.cfi_startproc
	pushq	%r13
	.cfi_def_cfa_offset 16
	.cfi_offset 13, -16
	pushq	%r12
	.cfi_def_cfa_offset 24
	.cfi_offset 12, -24
	movq	%rdi, %r13
	pushq	%rbp
	.cfi_def_cfa_offset 32
	.cfi_offset 6, -32
	pushq	%rbx
	.cfi_def_cfa_offset 40
	.cfi_offset 3, -40
	movq	%rsi, %r12
	leaq	PQconninfoOptions(%rip), %rbx
	leaq	.LC2(%rip), %rsi
	subq	$24, %rsp
	.cfi_def_cfa_offset 64
	jmp	.L281
	.p2align 4,,10
	.p2align 3
.L285:
	movq	56(%rbx), %rbp
	testq	%rbp, %rbp
	js	.L282
	movq	%r12, %rdi
	call	conninfo_getval
	testq	%rax, %rax
	je	.L282
	addq	%r13, %rbp
	movq	0(%rbp), %rdi
	testq	%rdi, %rdi
	je	.L283
	movq	%rax, 8(%rsp)
	call	free@PLT
	movq	8(%rsp), %rax
.L283:
	movq	%rax, %rdi
	call	strdup@PLT
	testq	%rax, %rax
	movq	%rax, 0(%rbp)
	je	.L293
.L282:
	movq	64(%rbx), %rsi
.L281:
	addq	$64, %rbx
	testq	%rsi, %rsi
	jne	.L285
	addq	$24, %rsp
	.cfi_remember_state
	.cfi_def_cfa_offset 40
	movl	$1, %eax
	popq	%rbx
	.cfi_def_cfa_offset 32
	popq	%rbp
	.cfi_def_cfa_offset 24
	popq	%r12
	.cfi_def_cfa_offset 16
	popq	%r13
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L293:
	.cfi_restore_state
	leaq	928(%r13), %rdi
	leaq	.LC3(%rip), %rsi
	call	printfPQExpBuffer@PLT
	addq	$24, %rsp
	.cfi_def_cfa_offset 40
	xorl	%eax, %eax
	popq	%rbx
	.cfi_def_cfa_offset 32
	popq	%rbp
	.cfi_def_cfa_offset 24
	popq	%r12
	.cfi_def_cfa_offset 16
	popq	%r13
	.cfi_def_cfa_offset 8
	ret
	.cfi_endproc
.LFE917:
	.size	fillPGconn, .-fillPGconn
	.p2align 4,,15
	.type	restoreErrorMessage, @function
restoreErrorMessage:
.LFB936:
	.cfi_startproc
	pushq	%r12
	.cfi_def_cfa_offset 16
	.cfi_offset 12, -16
	movq	%rdi, %r12
	pushq	%rbp
	.cfi_def_cfa_offset 24
	.cfi_offset 6, -24
	pushq	%rbx
	.cfi_def_cfa_offset 32
	.cfi_offset 3, -32
	movq	%rsi, %rbx
	movq	928(%rdi), %rsi
	leaq	928(%r12), %rbp
	movq	%rbx, %rdi
	call	appendPQExpBufferStr@PLT
	movq	%rbp, %rdi
	call	resetPQExpBuffer@PLT
	movq	(%rbx), %rsi
	movq	%rbp, %rdi
	call	appendPQExpBufferStr@PLT
	cmpq	$0, 16(%rbx)
	je	.L295
	cmpq	$0, 944(%r12)
	jne	.L296
.L295:
	leaq	.LC3(%rip), %rsi
	movq	%rbp, %rdi
	xorl	%eax, %eax
	call	printfPQExpBuffer@PLT
.L296:
	movq	%rbx, %rdi
	popq	%rbx
	.cfi_def_cfa_offset 24
	popq	%rbp
	.cfi_def_cfa_offset 16
	popq	%r12
	.cfi_def_cfa_offset 8
	jmp	termPQExpBuffer@PLT
	.cfi_endproc
.LFE936:
	.size	restoreErrorMessage, .-restoreErrorMessage
	.p2align 4,,15
	.type	parse_comma_separated_list, @function
parse_comma_separated_list:
.LFB920:
	.cfi_startproc
	pushq	%r13
	.cfi_def_cfa_offset 16
	.cfi_offset 13, -16
	pushq	%r12
	.cfi_def_cfa_offset 24
	.cfi_offset 12, -24
	pushq	%rbp
	.cfi_def_cfa_offset 32
	.cfi_offset 6, -32
	pushq	%rbx
	.cfi_def_cfa_offset 40
	.cfi_offset 3, -40
	movq	%rdi, %rbp
	subq	$8, %rsp
	.cfi_def_cfa_offset 48
	movq	(%rdi), %r13
	movzbl	0(%r13), %eax
	movq	%r13, %rbx
	testb	%al, %al
	je	.L303
	cmpb	$44, %al
	jne	.L300
	jmp	.L303
	.p2align 4,,10
	.p2align 3
.L319:
	cmpb	$44, %al
	je	.L304
.L300:
	addq	$1, %rbx
	movzbl	(%rbx), %eax
	testb	%al, %al
	jne	.L319
.L304:
	movq	%rbx, %rdi
	subq	%r13, %rdi
	movslq	%edi, %r12
	addl	$1, %edi
	movslq	%edi, %rdi
.L299:
	cmpb	$44, %al
	sete	(%rsi)
	call	malloc@PLT
	testq	%rax, %rax
	movq	%rax, %rcx
	je	.L302
	movq	%r12, %rdx
	movq	%r13, %rsi
	movq	%rax, %rdi
	call	memcpy@PLT
	movq	%rax, %rcx
	movb	$0, (%rax,%r12)
.L302:
	addq	$1, %rbx
	movq	%rcx, %rax
	movq	%rbx, 0(%rbp)
	addq	$8, %rsp
	.cfi_remember_state
	.cfi_def_cfa_offset 40
	popq	%rbx
	.cfi_def_cfa_offset 32
	popq	%rbp
	.cfi_def_cfa_offset 24
	popq	%r12
	.cfi_def_cfa_offset 16
	popq	%r13
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L303:
	.cfi_restore_state
	xorl	%r12d, %r12d
	movl	$1, %edi
	jmp	.L299
	.cfi_endproc
.LFE920:
	.size	parse_comma_separated_list, .-parse_comma_separated_list
	.p2align 4,,15
	.type	conninfo_init, @function
conninfo_init:
.LFB957:
	.cfi_startproc
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset 6, -16
	pushq	%rbx
	.cfi_def_cfa_offset 24
	.cfi_offset 3, -24
	movq	%rdi, %rbp
	movl	$2128, %edi
	subq	$8, %rsp
	.cfi_def_cfa_offset 32
	call	malloc@PLT
	testq	%rax, %rax
	movq	%rax, %rbx
	leaq	PQconninfoOptions(%rip), %rdx
	jne	.L321
	jmp	.L330
	.p2align 4,,10
	.p2align 3
.L325:
	movq	%rcx, %rax
.L321:
	movdqa	(%rdx), %xmm0
	addq	$64, %rdx
	movq	-16(%rdx), %rcx
	movups	%xmm0, (%rax)
	movdqa	-48(%rdx), %xmm0
	movq	%rcx, 48(%rax)
	leaq	56(%rax), %rcx
	movups	%xmm0, 16(%rax)
	movdqa	-32(%rdx), %xmm0
	cmpq	$0, (%rdx)
	movups	%xmm0, 32(%rax)
	jne	.L325
	addq	$112, %rax
	.p2align 4,,10
	.p2align 3
.L323:
	addq	$8, %rcx
	movq	$0, -8(%rcx)
	cmpq	%rcx, %rax
	ja	.L323
.L320:
	addq	$8, %rsp
	.cfi_remember_state
	.cfi_def_cfa_offset 24
	movq	%rbx, %rax
	popq	%rbx
	.cfi_def_cfa_offset 16
	popq	%rbp
	.cfi_def_cfa_offset 8
	ret
.L330:
	.cfi_restore_state
	leaq	.LC3(%rip), %rsi
	movq	%rbp, %rdi
	xorl	%eax, %eax
	call	printfPQExpBuffer@PLT
	jmp	.L320
	.cfi_endproc
.LFE957:
	.size	conninfo_init, .-conninfo_init
	.section	.rodata.str1.1
.LC4:
	.string	"r"
.LC5:
	.string	"service file \"%s\" not found\n"
	.section	.rodata.str1.8
	.align 8
.LC6:
	.string	"line %d too long in service file \"%s\"\n"
	.align 8
.LC7:
	.string	"syntax error in service file \"%s\", line %d\n"
	.align 8
.LC8:
	.string	"nested service specifications not supported in service file \"%s\", line %d\n"
	.text
	.p2align 4,,15
	.type	parseServiceFile, @function
parseServiceFile:
.LFB955:
	.cfi_startproc
	pushq	%r15
	.cfi_def_cfa_offset 16
	.cfi_offset 15, -16
	pushq	%r14
	.cfi_def_cfa_offset 24
	.cfi_offset 14, -24
	pushq	%r13
	.cfi_def_cfa_offset 32
	.cfi_offset 13, -32
	pushq	%r12
	.cfi_def_cfa_offset 40
	.cfi_offset 12, -40
	movq	%rdx, %r13
	pushq	%rbp
	.cfi_def_cfa_offset 48
	.cfi_offset 6, -48
	pushq	%rbx
	.cfi_def_cfa_offset 56
	.cfi_offset 3, -56
	xorl	%ebp, %ebp
	subq	$328, %rsp
	.cfi_def_cfa_offset 384
	movq	%rsi, 16(%rsp)
	leaq	.LC4(%rip), %rsi
	movq	%rdi, 40(%rsp)
	movq	%fs:40, %rax
	movq	%rax, 312(%rsp)
	xorl	%eax, %eax
	movq	%rcx, 32(%rsp)
	movq	%r8, 8(%rsp)
	call	fopen@PLT
	movq	%rax, %r12
	leaq	48(%rsp), %rax
	testq	%r12, %r12
	movq	%rax, (%rsp)
	je	.L372
	.p2align 4,,10
	.p2align 3
.L332:
	movq	(%rsp), %rdi
	movq	%r12, %rdx
	movl	$256, %esi
	call	fgets@PLT
	testq	%rax, %rax
	movq	%rax, %rbx
	je	.L370
	movq	%rbx, %rdi
	addl	$1, %ebp
	call	strlen@PLT
	cmpq	$254, %rax
	movq	%rax, %r15
	ja	.L373
	testq	%rax, %rax
	je	.L339
	call	__ctype_b_loc@PLT
	movslq	%r15d, %r9
	leal	-1(%r15), %ecx
	leal	-1(%r15), %edi
	leaq	-2(%rbx,%r9), %r8
	leaq	-1(%rbx,%r9), %rdx
	movslq	%edi, %rdi
	subq	%rcx, %r8
	jmp	.L338
	.p2align 4,,10
	.p2align 3
.L340:
	movq	%rdx, %rcx
	subq	$1, %rdx
	subq	%r9, %rcx
	cmpq	%rdx, %r8
	movb	$0, 1(%rcx,%rdi)
	je	.L339
.L338:
	movzbl	(%rdx), %esi
	movq	(%rax), %rcx
	testb	$32, 1(%rcx,%rsi,2)
	jne	.L340
.L339:
	movzbl	(%rbx), %r14d
	testb	%r14b, %r14b
	je	.L332
	call	__ctype_b_loc@PLT
	movq	(%rax), %rdx
	jmp	.L341
	.p2align 4,,10
	.p2align 3
.L342:
	addq	$1, %rbx
	movzbl	(%rbx), %r14d
	testb	%r14b, %r14b
	je	.L332
.L341:
	movzbl	%r14b, %eax
	testb	$32, 1(%rdx,%rax,2)
	jne	.L342
	cmpb	$35, %r14b
	je	.L332
	movq	8(%rsp), %rax
	cmpb	$91, %r14b
	movzbl	(%rax), %eax
	je	.L374
	testb	%al, %al
	je	.L332
	movl	$61, %esi
	movq	%rbx, %rdi
	call	strchr@PLT
	testq	%rax, %rax
	movq	%rax, 24(%rsp)
	je	.L375
	movq	24(%rsp), %rax
	leaq	.LC2(%rip), %rdi
	movl	$8, %ecx
	movq	%rbx, %rsi
	movb	$0, (%rax)
	repz cmpsb
	seta	%al
	sbbb	$0, %al
	testb	%al, %al
	je	.L347
	movq	0(%r13), %rdi
	testq	%rdi, %rdi
	je	.L349
	movq	%r13, %r14
	xorl	%r15d, %r15d
	jmp	.L348
	.p2align 4,,10
	.p2align 3
.L350:
	addl	$1, %r15d
	movslq	%r15d, %rdx
	leaq	0(,%rdx,8), %rax
	subq	%rdx, %rax
	leaq	0(%r13,%rax,8), %r14
	movq	(%r14), %rdi
	testq	%rdi, %rdi
	je	.L349
.L348:
	movq	%rbx, %rsi
	call	strcmp@PLT
	testl	%eax, %eax
	jne	.L350
	cmpq	$0, 24(%r14)
	jne	.L332
	movq	24(%rsp), %rdi
	addq	$1, %rdi
	call	strdup@PLT
	testq	%rax, %rax
	movq	%rax, 24(%r14)
	jne	.L332
	movq	32(%rsp), %rdi
	leaq	.LC3(%rip), %rsi
	call	printfPQExpBuffer@PLT
	movq	%r12, %rdi
	call	fclose@PLT
	movl	$3, %eax
	jmp	.L331
	.p2align 4,,10
	.p2align 3
.L374:
	testb	%al, %al
	jne	.L370
	movq	16(%rsp), %r15
	movq	%r15, %rdi
	call	strlen@PLT
	leaq	1(%rbx), %rdi
	movq	%rax, %rdx
	movq	%r15, %rsi
	movq	%rax, %r14
	call	strncmp@PLT
	testl	%eax, %eax
	jne	.L332
	cmpb	$93, 1(%rbx,%r14)
	jne	.L332
	movq	8(%rsp), %rax
	movb	$1, (%rax)
	jmp	.L332
	.p2align 4,,10
	.p2align 3
.L349:
	movq	40(%rsp), %rdx
	leaq	.LC7(%rip), %rsi
	movl	%ebp, %ecx
.L371:
	movq	32(%rsp), %rdi
	xorl	%eax, %eax
	call	printfPQExpBuffer@PLT
	movq	%r12, %rdi
	call	fclose@PLT
	movl	$3, %eax
.L331:
	movq	312(%rsp), %rsi
	xorq	%fs:40, %rsi
	jne	.L376
	addq	$328, %rsp
	.cfi_remember_state
	.cfi_def_cfa_offset 56
	popq	%rbx
	.cfi_def_cfa_offset 48
	popq	%rbp
	.cfi_def_cfa_offset 40
	popq	%r12
	.cfi_def_cfa_offset 32
	popq	%r13
	.cfi_def_cfa_offset 24
	popq	%r14
	.cfi_def_cfa_offset 16
	popq	%r15
	.cfi_def_cfa_offset 8
	ret
.L370:
	.cfi_restore_state
	movq	%r12, %rdi
	call	fclose@PLT
	xorl	%eax, %eax
	jmp	.L331
.L372:
	movq	40(%rsp), %rdx
	movq	32(%rsp), %rdi
	leaq	.LC5(%rip), %rsi
	xorl	%eax, %eax
	call	printfPQExpBuffer@PLT
	movl	$1, %eax
	jmp	.L331
.L373:
	movq	%r12, %rdi
	call	fclose@PLT
	movq	40(%rsp), %rcx
	movq	32(%rsp), %rdi
	leaq	.LC6(%rip), %rsi
	xorl	%eax, %eax
	movl	%ebp, %edx
	call	printfPQExpBuffer@PLT
	movl	$2, %eax
	jmp	.L331
.L347:
	movl	%ebp, %ecx
	movq	40(%rsp), %rdx
	leaq	.LC8(%rip), %rsi
	jmp	.L371
.L375:
	movq	40(%rsp), %rdx
	movq	32(%rsp), %rdi
	leaq	.LC7(%rip), %rsi
	movl	%ebp, %ecx
	call	printfPQExpBuffer@PLT
	movq	%r12, %rdi
	call	fclose@PLT
	movl	$3, %eax
	jmp	.L331
.L376:
	call	__stack_chk_fail@PLT
	.cfi_endproc
.LFE955:
	.size	parseServiceFile, .-parseServiceFile
	.section	.rodata.str1.8
	.align 8
.LC9:
	.string	"PQcancel() -- socket() failed: "
	.align 8
.LC10:
	.string	"PQcancel() -- connect() failed: "
	.section	.rodata.str1.1
.LC11:
	.string	"PQcancel() -- send() failed: "
	.text
	.p2align 4,,15
	.type	internal_cancel, @function
internal_cancel:
.LFB950:
	.cfi_startproc
	pushq	%r15
	.cfi_def_cfa_offset 16
	.cfi_offset 15, -16
	pushq	%r14
	.cfi_def_cfa_offset 24
	.cfi_offset 14, -24
	movq	%rdi, %r14
	pushq	%r13
	.cfi_def_cfa_offset 32
	.cfi_offset 13, -32
	pushq	%r12
	.cfi_def_cfa_offset 40
	.cfi_offset 12, -40
	movq	%rcx, %r13
	pushq	%rbp
	.cfi_def_cfa_offset 48
	.cfi_offset 6, -48
	pushq	%rbx
	.cfi_def_cfa_offset 56
	.cfi_offset 3, -56
	movl	%r8d, %r12d
	subq	$312, %rsp
	.cfi_def_cfa_offset 368
	movl	%esi, 8(%rsp)
	movl	%edx, 12(%rsp)
	movq	%fs:40, %rax
	movq	%rax, 296(%rsp)
	xorl	%eax, %eax
	call	__errno_location@PLT
	movzwl	(%r14), %edi
	xorl	%edx, %edx
	movl	$1, %esi
	movq	%rax, %rbx
	movl	(%rax), %r15d
	call	socket@PLT
	cmpl	$-1, %eax
	movl	%eax, %ebp
	jne	.L381
	jmp	.L399
	.p2align 4,,10
	.p2align 3
.L401:
	cmpl	$4, (%rbx)
	jne	.L400
.L381:
	movl	128(%r14), %edx
	movq	%r14, %rsi
	movl	%ebp, %edi
	call	connect@PLT
	testl	%eax, %eax
	js	.L401
	movabsq	$3321072690122457088, %rax
	movl	8(%rsp), %r14d
	movq	%rax, 16(%rsp)
	movl	12(%rsp), %eax
	bswap	%r14d
	bswap	%eax
	movl	%r14d, 24(%rsp)
	leaq	16(%rsp), %r14
	movl	%eax, 28(%rsp)
	jmp	.L384
	.p2align 4,,10
	.p2align 3
.L403:
	cmpl	$4, (%rbx)
	jne	.L402
.L384:
	xorl	%ecx, %ecx
	movl	$16, %edx
	movq	%r14, %rsi
	movl	%ebp, %edi
	call	send@PLT
	cmpq	$16, %rax
	jne	.L403
	jmp	.L385
	.p2align 4,,10
	.p2align 3
.L404:
	cmpl	$4, (%rbx)
	jne	.L386
.L385:
	xorl	%ecx, %ecx
	movl	$1, %edx
	movq	%r14, %rsi
	movl	%ebp, %edi
	call	recv@PLT
	testq	%rax, %rax
	js	.L404
.L386:
	movl	%ebp, %edi
	call	close@PLT
	movl	%r15d, (%rbx)
	movl	$1, %eax
.L377:
	movq	296(%rsp), %rcx
	xorq	%fs:40, %rcx
	jne	.L405
	addq	$312, %rsp
	.cfi_remember_state
	.cfi_def_cfa_offset 56
	popq	%rbx
	.cfi_def_cfa_offset 48
	popq	%rbp
	.cfi_def_cfa_offset 40
	popq	%r12
	.cfi_def_cfa_offset 32
	popq	%r13
	.cfi_def_cfa_offset 24
	popq	%r14
	.cfi_def_cfa_offset 16
	popq	%r15
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L400:
	.cfi_restore_state
	leaq	.LC10(%rip), %rsi
	movslq	%r12d, %rdx
	movq	%r13, %rdi
	call	strlcpy@PLT
.L383:
	movq	%r13, %rdi
	subl	$2, %r12d
	call	strlen@PLT
	subl	%eax, %r12d
	js	.L388
	movl	(%rbx), %edi
	leaq	32(%rsp), %rsi
	movl	$256, %edx
	call	pg_strerror_r@PLT
	movslq	%r12d, %rdx
	movq	%rax, %rsi
	movq	%r13, %rdi
	call	strncat@PLT
	movq	%r13, %rdi
	call	strlen@PLT
	movl	$10, %edx
	movw	%dx, 0(%r13,%rax)
.L388:
	movl	%ebp, %edi
	call	close@PLT
.L380:
	movl	%r15d, (%rbx)
	xorl	%eax, %eax
	jmp	.L377
	.p2align 4,,10
	.p2align 3
.L399:
	leaq	.LC9(%rip), %rsi
	movslq	%r12d, %rdx
	movq	%r13, %rdi
	subl	$2, %r12d
	call	strlcpy@PLT
	movq	%r13, %rdi
	call	strlen@PLT
	subl	%eax, %r12d
	js	.L380
	movl	(%rbx), %edi
	leaq	32(%rsp), %rsi
	movl	$256, %edx
	call	pg_strerror_r@PLT
	movq	%r13, %rdi
	movslq	%r12d, %rdx
	movq	%rax, %rsi
	call	strncat@PLT
	xorl	%eax, %eax
	orq	$-1, %rcx
	movq	%r13, %rdi
	repnz scasb
	movq	%rcx, %rax
	notq	%rax
	movw	$10, -1(%r13,%rax)
	jmp	.L380
	.p2align 4,,10
	.p2align 3
.L402:
	leaq	.LC11(%rip), %rsi
	movslq	%r12d, %rdx
	movq	%r13, %rdi
	call	strlcpy@PLT
	jmp	.L383
.L405:
	call	__stack_chk_fail@PLT
	.cfi_endproc
.LFE950:
	.size	internal_cancel, .-internal_cancel
	.p2align 4,,15
	.type	get_hexdigit, @function
get_hexdigit:
.LFB968:
	.cfi_startproc
	leal	-48(%rdi), %eax
	movsbl	%dil, %edx
	cmpb	$9, %al
	jbe	.L411
	leal	-65(%rdi), %eax
	cmpb	$5, %al
	jbe	.L412
	subl	$97, %edi
	xorl	%eax, %eax
	cmpb	$5, %dil
	ja	.L406
	subl	$87, %edx
	movl	$1, %eax
	movl	%edx, (%rsi)
.L406:
	rep ret
	.p2align 4,,10
	.p2align 3
.L412:
	subl	$55, %edx
	movl	$1, %eax
	movl	%edx, (%rsi)
	ret
	.p2align 4,,10
	.p2align 3
.L411:
	subl	$48, %edx
	movl	$1, %eax
	movl	%edx, (%rsi)
	ret
	.cfi_endproc
.LFE968:
	.size	get_hexdigit, .-get_hexdigit
	.section	.rodata.str1.8
	.align 8
.LC12:
	.string	"invalid percent-encoded token: \"%s\"\n"
	.align 8
.LC13:
	.string	"forbidden value %%00 in percent-encoded value: \"%s\"\n"
	.text
	.p2align 4,,15
	.type	conninfo_uri_decode, @function
conninfo_uri_decode:
.LFB967:
	.cfi_startproc
	pushq	%r12
	.cfi_def_cfa_offset 16
	.cfi_offset 12, -16
	pushq	%rbp
	.cfi_def_cfa_offset 24
	.cfi_offset 6, -24
	movq	%rsi, %r12
	pushq	%rbx
	.cfi_def_cfa_offset 32
	.cfi_offset 3, -32
	movq	%rdi, %rbp
	subq	$16, %rsp
	.cfi_def_cfa_offset 48
	movq	%fs:40, %rax
	movq	%rax, 8(%rsp)
	xorl	%eax, %eax
	call	strlen@PLT
	leaq	1(%rax), %rdi
	call	malloc@PLT
	testq	%rax, %rax
	movq	%rax, %rbx
	je	.L414
	leaq	1(%rax), %r8
	movq	%rbp, %rcx
	movq	%rsp, %r9
	leaq	4(%rsp), %r10
	jmp	.L415
	.p2align 4,,10
	.p2align 3
.L432:
	addq	$1, %rcx
	testb	%al, %al
	movb	%al, -1(%r8)
	je	.L413
.L418:
	addq	$1, %r8
.L415:
	movzbl	(%rcx), %eax
	cmpb	$37, %al
	jne	.L432
	movsbl	1(%rcx), %edi
	movq	%r9, %rsi
	call	get_hexdigit
	testb	%al, %al
	je	.L421
	movsbl	2(%rcx), %edi
	movq	%r10, %rsi
	leaq	3(%rcx), %r11
	call	get_hexdigit
	testb	%al, %al
	je	.L421
	movl	(%rsp), %eax
	sall	$4, %eax
	orl	4(%rsp), %eax
	je	.L433
	movb	%al, -1(%r8)
	movq	%r11, %rcx
	jmp	.L418
	.p2align 4,,10
	.p2align 3
.L421:
	leaq	.LC12(%rip), %rsi
	movq	%rbp, %rdx
.L431:
	movq	%r12, %rdi
	xorl	%eax, %eax
	call	printfPQExpBuffer@PLT
	movq	%rbx, %rdi
	xorl	%ebx, %ebx
	call	free@PLT
.L413:
	movq	8(%rsp), %rdx
	xorq	%fs:40, %rdx
	movq	%rbx, %rax
	jne	.L434
	addq	$16, %rsp
	.cfi_remember_state
	.cfi_def_cfa_offset 32
	popq	%rbx
	.cfi_def_cfa_offset 24
	popq	%rbp
	.cfi_def_cfa_offset 16
	popq	%r12
	.cfi_def_cfa_offset 8
	ret
.L433:
	.cfi_restore_state
	movq	%rbp, %rdx
	leaq	.LC13(%rip), %rsi
	jmp	.L431
.L414:
	leaq	.LC3(%rip), %rsi
	movq	%r12, %rdi
	xorl	%eax, %eax
	call	printfPQExpBuffer@PLT
	jmp	.L413
.L434:
	call	__stack_chk_fail@PLT
	.cfi_endproc
.LFE967:
	.size	conninfo_uri_decode, .-conninfo_uri_decode
	.section	.rodata.str1.1
.LC14:
	.string	"require"
.LC15:
	.string	"sslmode"
.LC16:
	.string	"prefer"
.LC17:
	.string	"requiressl"
	.section	.rodata.str1.8
	.align 8
.LC18:
	.string	"invalid connection option \"%s\"\n"
	.text
	.p2align 4,,15
	.type	conninfo_storeval, @function
conninfo_storeval:
.LFB970:
	.cfi_startproc
	pushq	%r15
	.cfi_def_cfa_offset 16
	.cfi_offset 15, -16
	pushq	%r14
	.cfi_def_cfa_offset 24
	.cfi_offset 14, -24
	movq	%rcx, %r15
	pushq	%r13
	.cfi_def_cfa_offset 32
	.cfi_offset 13, -32
	pushq	%r12
	.cfi_def_cfa_offset 40
	.cfi_offset 12, -40
	movl	$11, %ecx
	pushq	%rbp
	.cfi_def_cfa_offset 48
	.cfi_offset 6, -48
	pushq	%rbx
	.cfi_def_cfa_offset 56
	.cfi_offset 3, -56
	movq	%rdi, %rbx
	leaq	.LC17(%rip), %rdi
	movq	%rsi, %rbp
	movq	%rdx, %r12
	subq	$8, %rsp
	.cfi_def_cfa_offset 64
	movl	%r8d, %r14d
	movl	%r9d, %r13d
	repz cmpsb
	seta	%al
	sbbb	$0, %al
	testb	%al, %al
	jne	.L436
	cmpb	$49, (%rdx)
	je	.L447
	leaq	.LC16(%rip), %r12
	leaq	.LC15(%rip), %rbp
.L436:
	movq	(%rbx), %rdi
	testq	%rdi, %rdi
	jne	.L439
	jmp	.L437
	.p2align 4,,10
	.p2align 3
.L459:
	addq	$56, %rbx
	movq	(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L437
.L439:
	movq	%rbp, %rsi
	call	strcmp@PLT
	testl	%eax, %eax
	jne	.L459
	testb	%r13b, %r13b
	jne	.L445
	movq	%r12, %rdi
	call	strdup@PLT
	testq	%rax, %rax
	movq	%rax, %rbp
	je	.L460
.L441:
	movq	24(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L442
	call	free@PLT
.L442:
	movq	%rbp, 24(%rbx)
	movq	%rbx, %rax
.L435:
	addq	$8, %rsp
	.cfi_remember_state
	.cfi_def_cfa_offset 56
	popq	%rbx
	.cfi_def_cfa_offset 48
	popq	%rbp
	.cfi_def_cfa_offset 40
	popq	%r12
	.cfi_def_cfa_offset 32
	popq	%r13
	.cfi_def_cfa_offset 24
	popq	%r14
	.cfi_def_cfa_offset 16
	popq	%r15
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L437:
	.cfi_restore_state
	testb	%r14b, %r14b
	je	.L461
	xorl	%eax, %eax
.L462:
	addq	$8, %rsp
	.cfi_remember_state
	.cfi_def_cfa_offset 56
	popq	%rbx
	.cfi_def_cfa_offset 48
	popq	%rbp
	.cfi_def_cfa_offset 40
	popq	%r12
	.cfi_def_cfa_offset 32
	popq	%r13
	.cfi_def_cfa_offset 24
	popq	%r14
	.cfi_def_cfa_offset 16
	popq	%r15
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L445:
	.cfi_restore_state
	movq	%r15, %rsi
	movq	%r12, %rdi
	call	conninfo_uri_decode
	testq	%rax, %rax
	movq	%rax, %rbp
	jne	.L441
	xorl	%eax, %eax
	jmp	.L462
	.p2align 4,,10
	.p2align 3
.L461:
	leaq	.LC18(%rip), %rsi
	xorl	%eax, %eax
	movq	%rbp, %rdx
	movq	%r15, %rdi
	call	printfPQExpBuffer@PLT
	xorl	%eax, %eax
	jmp	.L435
	.p2align 4,,10
	.p2align 3
.L447:
	leaq	.LC14(%rip), %r12
	leaq	.LC15(%rip), %rbp
	jmp	.L436
.L460:
	leaq	.LC3(%rip), %rsi
	xorl	%eax, %eax
	movq	%r15, %rdi
	call	printfPQExpBuffer@PLT
	xorl	%eax, %eax
	jmp	.L462
	.cfi_endproc
.LFE970:
	.size	conninfo_storeval, .-conninfo_storeval
	.section	.rodata.str1.1
.LC19:
	.string	"failed to lock mutex"
.LC20:
	.string	"%s\n"
.LC21:
	.string	"failed to unlock mutex"
	.text
	.p2align 4,,15
	.type	default_threadlock, @function
default_threadlock:
.LFB1008:
	.cfi_startproc
	subq	$8, %rsp
	.cfi_def_cfa_offset 16
	testl	%edi, %edi
	leaq	singlethread_lock.25935(%rip), %rdi
	je	.L464
	call	pthread_mutex_lock@PLT
	testl	%eax, %eax
	jne	.L474
.L463:
	addq	$8, %rsp
	.cfi_remember_state
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L464:
	.cfi_restore_state
	call	pthread_mutex_unlock@PLT
	testl	%eax, %eax
	je	.L463
	leaq	.LC21(%rip), %rdx
.L473:
	movq	stderr@GOTPCREL(%rip), %rax
	leaq	.LC20(%rip), %rsi
	movq	(%rax), %rdi
	xorl	%eax, %eax
	call	pg_fprintf@PLT
	call	abort@PLT
.L474:
	leaq	.LC19(%rip), %rdx
	jmp	.L473
	.cfi_endproc
.LFE1008:
	.size	default_threadlock, .-default_threadlock
	.p2align 4,,15
	.type	uri_prefix_length, @function
uri_prefix_length:
.LFB959:
	.cfi_startproc
	pushq	%rbx
	.cfi_def_cfa_offset 16
	.cfi_offset 3, -16
	leaq	uri_designator(%rip), %rsi
	movl	$13, %edx
	movq	%rdi, %rbx
	call	strncmp@PLT
	testl	%eax, %eax
	movl	$13, %edx
	je	.L475
	leaq	short_uri_designator(%rip), %rsi
	movl	$11, %edx
	movq	%rbx, %rdi
	call	strncmp@PLT
	cmpl	$1, %eax
	sbbl	%edx, %edx
	andl	$11, %edx
.L475:
	movl	%edx, %eax
	popq	%rbx
	.cfi_def_cfa_offset 8
	ret
	.cfi_endproc
.LFE959:
	.size	uri_prefix_length, .-uri_prefix_length
	.p2align 4,,15
	.type	recognized_connection_string, @function
recognized_connection_string:
.LFB960:
	.cfi_startproc
	pushq	%rbx
	.cfi_def_cfa_offset 16
	.cfi_offset 3, -16
	movq	%rdi, %rbx
	call	uri_prefix_length
	testl	%eax, %eax
	movl	$1, %edx
	jne	.L482
	movl	$61, %esi
	movq	%rbx, %rdi
	call	strchr@PLT
	testq	%rax, %rax
	setne	%dl
.L482:
	movl	%edx, %eax
	popq	%rbx
	.cfi_def_cfa_offset 8
	ret
	.cfi_endproc
.LFE960:
	.size	recognized_connection_string, .-recognized_connection_string
	.section	.rodata.str1.1
.LC22:
	.string	"TLSv1"
.LC23:
	.string	"TLSv1.1"
.LC24:
	.string	"TLSv1.2"
.LC25:
	.string	"TLSv1.3"
	.text
	.p2align 4,,15
	.type	sslVerifyProtocolVersion, @function
sslVerifyProtocolVersion:
.LFB1005:
	.cfi_startproc
	testq	%rdi, %rdi
	je	.L501
	cmpb	$0, (%rdi)
	jne	.L504
.L501:
	movl	$1, %eax
	ret
	.p2align 4,,10
	.p2align 3
.L504:
	pushq	%rbx
	.cfi_def_cfa_offset 16
	.cfi_offset 3, -16
	leaq	.LC22(%rip), %rsi
	movq	%rdi, %rbx
	call	pg_strcasecmp@PLT
	testl	%eax, %eax
	jne	.L505
.L490:
	movl	$1, %eax
	popq	%rbx
	.cfi_remember_state
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L505:
	.cfi_restore_state
	leaq	.LC23(%rip), %rsi
	movq	%rbx, %rdi
	call	pg_strcasecmp@PLT
	testl	%eax, %eax
	je	.L490
	leaq	.LC24(%rip), %rsi
	movq	%rbx, %rdi
	call	pg_strcasecmp@PLT
	testl	%eax, %eax
	je	.L490
	leaq	.LC25(%rip), %rsi
	movq	%rbx, %rdi
	call	pg_strcasecmp@PLT
	testl	%eax, %eax
	sete	%al
	popq	%rbx
	.cfi_def_cfa_offset 8
	ret
	.cfi_endproc
.LFE1005:
	.size	sslVerifyProtocolVersion, .-sslVerifyProtocolVersion
	.p2align 4,,15
	.type	sendTerminateConn.part.7, @function
sendTerminateConn.part.7:
.LFB1017:
	.cfi_startproc
	pushq	%rbx
	.cfi_def_cfa_offset 16
	.cfi_offset 3, -16
	movq	%rdi, %rdx
	movq	%rdi, %rbx
	xorl	%esi, %esi
	movl	$88, %edi
	call	pqPutMsgStart@PLT
	movq	%rbx, %rdi
	call	pqPutMsgEnd@PLT
	movq	%rbx, %rdi
	popq	%rbx
	.cfi_def_cfa_offset 8
	jmp	pqFlush@PLT
	.cfi_endproc
.LFE1017:
	.size	sendTerminateConn.part.7, .-sendTerminateConn.part.7
	.p2align 4,,15
	.type	getHostaddr.constprop.14, @function
getHostaddr.constprop.14:
.LFB1024:
	.cfi_startproc
	pushq	%rbx
	.cfi_def_cfa_offset 16
	.cfi_offset 3, -16
	movzwl	560(%rdi), %eax
	movq	%rsi, %rbx
	cmpw	$2, %ax
	je	.L521
	cmpw	$10, %ax
	je	.L522
.L510:
	movb	$0, (%rbx)
	popq	%rbx
	.cfi_remember_state
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L521:
	.cfi_restore_state
	leaq	564(%rdi), %rsi
	movl	$1025, %r8d
	movq	%rbx, %rcx
	movl	$32, %edx
	movl	$2, %edi
	call	pg_inet_net_ntop@PLT
	testq	%rax, %rax
	je	.L510
	popq	%rbx
	.cfi_remember_state
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L522:
	.cfi_restore_state
	leaq	568(%rdi), %rsi
	movl	$1025, %r8d
	movq	%rbx, %rcx
	movl	$128, %edx
	movl	$10, %edi
	call	pg_inet_net_ntop@PLT
	testq	%rax, %rax
	je	.L510
	popq	%rbx
	.cfi_def_cfa_offset 8
	ret
	.cfi_endproc
.LFE1024:
	.size	getHostaddr.constprop.14, .-getHostaddr.constprop.14
	.section	.rodata.str1.1
.LC26:
	.string	"5432"
	.section	.rodata.str1.8
	.align 8
.LC27:
	.string	"could not connect to server: %s\n\tIs the server running locally and accepting\n\tconnections on Unix domain socket \"%s\"?\n"
	.align 8
.LC28:
	.string	"could not connect to server: %s\n\tIs the server running on host \"%s\" (%s) and accepting\n\tTCP/IP connections on port %s?\n"
	.align 8
.LC29:
	.string	"could not connect to server: %s\n\tIs the server running on host \"%s\" and accepting\n\tTCP/IP connections on port %s?\n"
	.text
	.p2align 4,,15
	.type	connectFailureMessage, @function
connectFailureMessage:
.LFB926:
	.cfi_startproc
	pushq	%r14
	.cfi_def_cfa_offset 16
	.cfi_offset 14, -16
	pushq	%r13
	.cfi_def_cfa_offset 24
	.cfi_offset 13, -24
	pushq	%r12
	.cfi_def_cfa_offset 32
	.cfi_offset 12, -32
	pushq	%rbp
	.cfi_def_cfa_offset 40
	.cfi_offset 6, -40
	movl	%esi, %r12d
	pushq	%rbx
	.cfi_def_cfa_offset 48
	.cfi_offset 3, -48
	leaq	928(%rdi), %rbp
	movq	%rdi, %rbx
	subq	$1296, %rsp
	.cfi_def_cfa_offset 1344
	movq	%fs:40, %rax
	movq	%rax, 1288(%rsp)
	xorl	%eax, %eax
	cmpw	$1, 560(%rdi)
	leaq	256(%rsp), %r13
	je	.L542
	movq	%r13, %rsi
	call	getHostaddr.constprop.14
	movslq	396(%rbx), %rax
	leaq	(%rax,%rax,4), %rdx
	movq	400(%rbx), %rax
	leaq	(%rax,%rdx,8), %rax
	cmpl	$1, (%rax)
	je	.L543
	movq	24(%rax), %rbx
	movq	8(%rax), %r14
	testq	%rbx, %rbx
	je	.L534
	cmpb	$0, (%rbx)
	leaq	.LC26(%rip), %rax
	cmove	%rax, %rbx
.L531:
	cmpb	$0, 256(%rsp)
	je	.L528
	movq	%r13, %rsi
	movq	%r14, %rdi
	call	strcmp@PLT
	testl	%eax, %eax
	je	.L528
	movq	%rsp, %rsi
	movl	%r12d, %edi
	movl	$256, %edx
	call	pg_strerror_r@PLT
	leaq	.LC28(%rip), %rsi
	movq	%rax, %rdx
	movq	%rbx, %r9
	movq	%r13, %r8
	movq	%r14, %rcx
	movq	%rbp, %rdi
	xorl	%eax, %eax
	call	appendPQExpBuffer@PLT
	jmp	.L523
	.p2align 4,,10
	.p2align 3
.L543:
	movq	24(%rax), %rbx
	movq	16(%rax), %r14
	testq	%rbx, %rbx
	je	.L544
	cmpb	$0, (%rbx)
	leaq	.LC26(%rip), %rax
	cmove	%rax, %rbx
.L528:
	movq	%rsp, %rsi
	movl	$256, %edx
	movl	%r12d, %edi
	call	pg_strerror_r@PLT
	leaq	.LC29(%rip), %rsi
	movq	%rax, %rdx
	movq	%rbx, %r8
	movq	%r14, %rcx
	movq	%rbp, %rdi
	xorl	%eax, %eax
	call	appendPQExpBuffer@PLT
.L523:
	movq	1288(%rsp), %rax
	xorq	%fs:40, %rax
	jne	.L545
	addq	$1296, %rsp
	.cfi_remember_state
	.cfi_def_cfa_offset 48
	popq	%rbx
	.cfi_def_cfa_offset 40
	popq	%rbp
	.cfi_def_cfa_offset 32
	popq	%r12
	.cfi_def_cfa_offset 24
	popq	%r13
	.cfi_def_cfa_offset 16
	popq	%r14
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L542:
	.cfi_restore_state
	subq	$8, %rsp
	.cfi_def_cfa_offset 1352
	movl	688(%rdi), %esi
	leaq	560(%rdi), %rdi
	pushq	$2
	.cfi_def_cfa_offset 1360
	movl	$1025, %r9d
	movq	%r13, %r8
	xorl	%ecx, %ecx
	xorl	%edx, %edx
	call	pg_getnameinfo_all@PLT
	leaq	16(%rsp), %rsi
	movl	$256, %edx
	movl	%r12d, %edi
	call	pg_strerror_r@PLT
	leaq	.LC27(%rip), %rsi
	movq	%rax, %rdx
	movq	%r13, %rcx
	xorl	%eax, %eax
	movq	%rbp, %rdi
	call	appendPQExpBuffer@PLT
	popq	%rax
	.cfi_def_cfa_offset 1352
	popq	%rdx
	.cfi_def_cfa_offset 1344
	jmp	.L523
	.p2align 4,,10
	.p2align 3
.L534:
	leaq	.LC26(%rip), %rbx
	jmp	.L531
.L545:
	call	__stack_chk_fail@PLT
.L544:
	leaq	.LC26(%rip), %rbx
	jmp	.L528
	.cfi_endproc
.LFE926:
	.size	connectFailureMessage, .-connectFailureMessage
	.p2align 4,,15
	.globl	pqDropConnection
	.type	pqDropConnection, @function
pqDropConnection:
.LFB909:
	.cfi_startproc
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset 6, -16
	pushq	%rbx
	.cfi_def_cfa_offset 24
	.cfi_offset 3, -24
	movq	%rdi, %rbx
	movl	%esi, %ebp
	subq	$8, %rsp
	.cfi_def_cfa_offset 32
	call	pqsecure_close@PLT
	movl	416(%rbx), %edi
	cmpl	$-1, %edi
	je	.L547
	call	close@PLT
.L547:
	testb	%bpl, %bpl
	movl	$-1, 416(%rbx)
	je	.L548
	movq	$0, 820(%rbx)
	movl	$0, 828(%rbx)
.L548:
	movq	888(%rbx), %rdi
	movl	$0, 844(%rbx)
	testq	%rdi, %rdi
	je	.L546
	call	pg_fe_scram_free@PLT
	movq	$0, 888(%rbx)
.L546:
	addq	$8, %rsp
	.cfi_def_cfa_offset 24
	popq	%rbx
	.cfi_def_cfa_offset 16
	popq	%rbp
	.cfi_def_cfa_offset 8
	ret
	.cfi_endproc
.LFE909:
	.size	pqDropConnection, .-pqDropConnection
	.p2align 4,,15
	.type	closePGconn, @function
closePGconn:
.LFB943:
	.cfi_startproc
	pushq	%rbx
	.cfi_def_cfa_offset 16
	.cfi_offset 3, -16
	cmpl	$-1, 416(%rdi)
	movq	%rdi, %rbx
	je	.L558
	movl	336(%rdi), %eax
	testl	%eax, %eax
	je	.L564
.L558:
	movl	$1, %esi
	movq	%rbx, %rdi
	movb	$0, 367(%rbx)
	call	pqDropConnection@PLT
	movq	%rbx, %rdi
	movl	$1, 336(%rbx)
	movl	$0, 340(%rbx)
	movl	$0, 344(%rbx)
	call	pqClearAsyncResult@PLT
	leaq	928(%rbx), %rdi
	call	resetPQExpBuffer@PLT
	movq	728(%rbx), %rsi
	testq	%rsi, %rsi
	je	.L559
	movl	744(%rbx), %edi
	call	pg_freeaddrinfo_all@PLT
	movq	$0, 728(%rbx)
	movq	$0, 736(%rbx)
.L559:
	movq	%rbx, %rdi
	popq	%rbx
	.cfi_remember_state
	.cfi_def_cfa_offset 8
	jmp	pqDropServerData
	.p2align 4,,10
	.p2align 3
.L564:
	.cfi_restore_state
	call	sendTerminateConn.part.7
	jmp	.L558
	.cfi_endproc
.LFE943:
	.size	closePGconn, .-closePGconn
	.p2align 4,,15
	.globl	PQfinish
	.type	PQfinish, @function
PQfinish:
.LFB944:
	.cfi_startproc
	testq	%rdi, %rdi
	je	.L565
	pushq	%rbx
	.cfi_def_cfa_offset 16
	.cfi_offset 3, -16
	movq	%rdi, %rbx
	call	closePGconn
	movq	%rbx, %rdi
	popq	%rbx
	.cfi_restore 3
	.cfi_def_cfa_offset 8
	jmp	freePGconn
	.p2align 4,,10
	.p2align 3
.L565:
	rep ret
	.cfi_endproc
.LFE944:
	.size	PQfinish, .-PQfinish
	.p2align 4,,15
	.globl	PQgetCancel
	.type	PQgetCancel, @function
PQgetCancel:
.LFB948:
	.cfi_startproc
	testq	%rdi, %rdi
	je	.L573
	cmpl	$-1, 416(%rdi)
	je	.L573
	pushq	%rbx
	.cfi_def_cfa_offset 16
	.cfi_offset 3, -16
	movq	%rdi, %rbx
	movl	$144, %edi
	call	malloc@PLT
	testq	%rax, %rax
	je	.L570
	movdqu	560(%rbx), %xmm0
	movq	688(%rbx), %rdx
	movups	%xmm0, (%rax)
	movdqu	576(%rbx), %xmm0
	movq	%rdx, 128(%rax)
	movl	764(%rbx), %edx
	movups	%xmm0, 16(%rax)
	movdqu	592(%rbx), %xmm0
	movl	%edx, 136(%rax)
	movl	768(%rbx), %edx
	movups	%xmm0, 32(%rax)
	movdqu	608(%rbx), %xmm0
	movl	%edx, 140(%rax)
	movups	%xmm0, 48(%rax)
	movdqu	624(%rbx), %xmm0
	movups	%xmm0, 64(%rax)
	movdqu	640(%rbx), %xmm0
	movups	%xmm0, 80(%rax)
	movdqu	656(%rbx), %xmm0
	movups	%xmm0, 96(%rax)
	movdqu	672(%rbx), %xmm0
	movups	%xmm0, 112(%rax)
.L570:
	popq	%rbx
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L573:
	.cfi_restore 3
	xorl	%eax, %eax
	ret
	.cfi_endproc
.LFE948:
	.size	PQgetCancel, .-PQgetCancel
	.p2align 4,,15
	.globl	PQfreeCancel
	.type	PQfreeCancel, @function
PQfreeCancel:
.LFB949:
	.cfi_startproc
	testq	%rdi, %rdi
	je	.L581
	jmp	free@PLT
	.p2align 4,,10
	.p2align 3
.L581:
	rep ret
	.cfi_endproc
.LFE949:
	.size	PQfreeCancel, .-PQfreeCancel
	.section	.rodata.str1.8
	.align 8
.LC30:
	.string	"PQcancel() -- no cancel object supplied"
	.text
	.p2align 4,,15
	.globl	PQcancel
	.type	PQcancel, @function
PQcancel:
.LFB951:
	.cfi_startproc
	testq	%rdi, %rdi
	je	.L587
	movl	%edx, %r8d
	movq	%rsi, %rcx
	movl	140(%rdi), %edx
	movl	136(%rdi), %esi
	jmp	internal_cancel
	.p2align 4,,10
	.p2align 3
.L587:
	movq	%rsi, %rdi
	leaq	.LC30(%rip), %rsi
	subq	$8, %rsp
	.cfi_def_cfa_offset 16
	movslq	%edx, %rdx
	call	strlcpy@PLT
	xorl	%eax, %eax
	addq	$8, %rsp
	.cfi_def_cfa_offset 8
	ret
	.cfi_endproc
.LFE951:
	.size	PQcancel, .-PQcancel
	.section	.rodata.str1.8
	.align 8
.LC31:
	.string	"PQrequestCancel() -- connection is not open\n"
	.text
	.p2align 4,,15
	.globl	PQrequestCancel
	.type	PQrequestCancel, @function
PQrequestCancel:
.LFB952:
	.cfi_startproc
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset 6, -16
	pushq	%rbx
	.cfi_def_cfa_offset 24
	.cfi_offset 3, -24
	xorl	%ebp, %ebp
	subq	$8, %rsp
	.cfi_def_cfa_offset 32
	testq	%rdi, %rdi
	je	.L590
	cmpl	$-1, 416(%rdi)
	movq	%rdi, %rbx
	movq	944(%rdi), %r8
	movq	928(%rdi), %rcx
	je	.L596
	movl	768(%rdi), %edx
	movl	764(%rdi), %esi
	leaq	560(%rdi), %rdi
	call	internal_cancel
	testl	%eax, %eax
	movl	%eax, %ebp
	jne	.L590
	movq	928(%rbx), %rdi
	call	strlen@PLT
	movq	%rax, 936(%rbx)
.L590:
	addq	$8, %rsp
	.cfi_remember_state
	.cfi_def_cfa_offset 24
	movl	%ebp, %eax
	popq	%rbx
	.cfi_def_cfa_offset 16
	popq	%rbp
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L596:
	.cfi_restore_state
	leaq	.LC31(%rip), %rsi
	movq	%rcx, %rdi
	movq	%r8, %rdx
	call	strlcpy@PLT
	movq	928(%rbx), %rdi
	call	strlen@PLT
	movq	%rax, 936(%rbx)
	addq	$8, %rsp
	.cfi_def_cfa_offset 24
	movl	%ebp, %eax
	popq	%rbx
	.cfi_def_cfa_offset 16
	popq	%rbp
	.cfi_def_cfa_offset 8
	ret
	.cfi_endproc
.LFE952:
	.size	PQrequestCancel, .-PQrequestCancel
	.p2align 4,,15
	.globl	pqPacketSend
	.type	pqPacketSend, @function
pqPacketSend:
.LFB953:
	.cfi_startproc
	pushq	%r12
	.cfi_def_cfa_offset 16
	.cfi_offset 12, -16
	pushq	%rbp
	.cfi_def_cfa_offset 24
	.cfi_offset 6, -24
	movq	%rdx, %rbp
	pushq	%rbx
	.cfi_def_cfa_offset 32
	.cfi_offset 3, -32
	movq	%rdi, %rbx
	movsbl	%sil, %edi
	movq	%rbx, %rdx
	movl	$1, %esi
	movq	%rcx, %r12
	call	pqPutMsgStart@PLT
	testl	%eax, %eax
	jne	.L600
	movq	%rbx, %rdx
	movq	%r12, %rsi
	movq	%rbp, %rdi
	call	pqPutnchar@PLT
	testl	%eax, %eax
	jne	.L600
	movq	%rbx, %rdi
	call	pqPutMsgEnd@PLT
	testl	%eax, %eax
	jne	.L600
	movq	%rbx, %rdi
	call	pqFlush@PLT
	testl	%eax, %eax
	setne	%al
	movzbl	%al, %eax
	popq	%rbx
	.cfi_remember_state
	.cfi_def_cfa_offset 24
	negl	%eax
	popq	%rbp
	.cfi_def_cfa_offset 16
	popq	%r12
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L600:
	.cfi_restore_state
	popq	%rbx
	.cfi_def_cfa_offset 24
	movl	$-1, %eax
	popq	%rbp
	.cfi_def_cfa_offset 16
	popq	%r12
	.cfi_def_cfa_offset 8
	ret
	.cfi_endproc
.LFE953:
	.size	pqPacketSend, .-pqPacketSend
	.section	.rodata.str1.8
	.align 8
.LC32:
	.string	"invalid connection state, probably indicative of memory corruption\n"
	.section	.rodata.str1.1
.LC33:
	.string	"port"
.LC34:
	.string	"invalid port number: \"%s\"\n"
.LC35:
	.string	"%d"
	.section	.rodata.str1.8
	.align 8
.LC36:
	.string	"could not translate host name \"%s\" to address: %s\n"
	.align 8
.LC37:
	.string	"could not parse network address \"%s\": %s\n"
	.section	.rodata.str1.1
.LC38:
	.string	"%s/.s.PGSQL.%d"
	.section	.rodata.str1.8
	.align 8
.LC39:
	.string	"Unix-domain socket path \"%s\" is too long (maximum %d bytes)\n"
	.align 8
.LC40:
	.string	"could not translate Unix-domain socket path \"%s\" to address: %s\n"
	.section	.rodata.str1.1
.LC41:
	.string	"could not create socket: %s\n"
	.section	.rodata.str1.8
	.align 8
.LC42:
	.string	"could not set socket to TCP no delay mode: %s\n"
	.align 8
.LC43:
	.string	"could not set socket to nonblocking mode: %s\n"
	.align 8
.LC44:
	.string	"could not set socket to close-on-exec mode: %s\n"
	.align 8
.LC45:
	.string	"keepalives parameter must be an integer\n"
	.section	.rodata.str1.1
.LC46:
	.string	"SO_KEEPALIVE"
.LC47:
	.string	"setsockopt(%s) failed: %s\n"
.LC48:
	.string	"keepalives_idle"
.LC49:
	.string	"TCP_KEEPIDLE"
.LC50:
	.string	"keepalives_interval"
.LC51:
	.string	"TCP_KEEPINTVL"
.LC52:
	.string	"keepalives_count"
.LC53:
	.string	"TCP_KEEPCNT"
.LC54:
	.string	"tcp_user_timeout"
.LC55:
	.string	"TCP_USER_TIMEOUT"
	.section	.rodata.str1.8
	.align 8
.LC56:
	.string	"could not get socket error status: %s\n"
	.align 8
.LC57:
	.string	"could not get client address from socket: %s\n"
	.align 8
.LC58:
	.string	"requirepeer parameter is not supported on this platform\n"
	.align 8
.LC59:
	.string	"could not get peer credentials: %s\n"
	.align 8
.LC60:
	.string	"could not look up local user ID %d: %s\n"
	.align 8
.LC61:
	.string	"local user with ID %d does not exist\n"
	.align 8
.LC62:
	.string	"requirepeer specifies \"%s\", but actual peer user name is \"%s\"\n"
	.align 8
.LC63:
	.string	"could not send SSL negotiation packet: %s\n"
	.align 8
.LC64:
	.string	"could not send startup packet: %s\n"
	.align 8
.LC65:
	.string	"server does not support SSL, but SSL was required\n"
	.align 8
.LC66:
	.string	"received invalid response to SSL negotiation: %c\n"
	.align 8
.LC67:
	.string	"expected authentication request from server, but received %c\n"
	.section	.rodata.str1.1
.LC68:
	.string	"28P01"
	.section	.rodata.str1.8
	.align 8
.LC69:
	.string	"password retrieved from file \"%s\"\n"
	.align 8
.LC70:
	.string	"unexpected message from server during startup\n"
	.section	.rodata.str1.1
.LC71:
	.string	"42704"
.LC72:
	.string	"read-write"
.LC73:
	.string	"SHOW transaction_read_only"
.LC74:
	.string	"on"
	.section	.rodata.str1.8
	.align 8
.LC75:
	.string	"could not make a writable connection to server \"%s:%s\"\n"
	.align 8
.LC76:
	.string	"test \"SHOW transaction_read_only\" failed on server \"%s:%s\"\n"
	.align 8
.LC77:
	.string	"invalid connection state %d, probably indicative of memory corruption\n"
	.text
	.p2align 4,,15
	.globl	PQconnectPoll
	.type	PQconnectPoll, @function
PQconnectPoll:
.LFB937:
	.cfi_startproc
	pushq	%r15
	.cfi_def_cfa_offset 16
	.cfi_offset 15, -16
	pushq	%r14
	.cfi_def_cfa_offset 24
	.cfi_offset 14, -24
	pushq	%r13
	.cfi_def_cfa_offset 32
	.cfi_offset 13, -32
	pushq	%r12
	.cfi_def_cfa_offset 40
	.cfi_offset 12, -40
	pushq	%rbp
	.cfi_def_cfa_offset 48
	.cfi_offset 6, -48
	pushq	%rbx
	.cfi_def_cfa_offset 56
	.cfi_offset 3, -56
	subq	$8856, %rsp
	.cfi_def_cfa_offset 8912
	movq	%fs:40, %rax
	movq	%rax, 8840(%rsp)
	xorl	%eax, %eax
	testq	%rdi, %rdi
	je	.L602
	cmpl	$11, 336(%rdi)
	movq	%rdi, %r12
	ja	.L604
	movl	336(%rdi), %eax
	leaq	.L605(%rip), %rdx
	movslq	(%rdx,%rax,4), %rax
	addq	%rdx, %rax
	jmp	*%rax
	.section	.rodata
	.align 4
	.align 4
.L605:
	.long	.L793-.L605
	.long	.L916-.L605
	.long	.L607-.L605
	.long	.L607-.L605
	.long	.L608-.L605
	.long	.L608-.L605
	.long	.L607-.L605
	.long	.L607-.L605
	.long	.L607-.L605
	.long	.L607-.L605
	.long	.L607-.L605
	.long	.L607-.L605
	.text
.L649:
	movl	$0, 336(%r12)
	movq	%r12, %rdi
	call	pqSetenvPoll@PLT
	cmpl	$2, %eax
	je	.L769
	cmpl	$3, %eax
	je	.L770
	cmpl	$1, %eax
	jne	.L611
	movl	$6, 336(%r12)
	.p2align 4,,10
	.p2align 3
.L602:
	movq	8840(%rsp), %rbx
	xorq	%fs:40, %rbx
	jne	.L928
	addq	$8856, %rsp
	.cfi_remember_state
	.cfi_def_cfa_offset 56
	popq	%rbx
	.cfi_def_cfa_offset 48
	popq	%rbp
	.cfi_def_cfa_offset 40
	popq	%r12
	.cfi_def_cfa_offset 32
	popq	%r13
	.cfi_def_cfa_offset 24
	popq	%r14
	.cfi_def_cfa_offset 16
	popq	%r15
	.cfi_def_cfa_offset 8
	ret
.L646:
	.cfi_restore_state
	movq	224(%r12), %rax
	movzwl	560(%r12), %edx
	leaq	80(%rsp), %r10
	testq	%rax, %rax
	je	.L698
	cmpb	$0, (%rax)
	je	.L698
	cmpw	$1, %dx
	je	.L929
.L699:
	cmpb	$0, 897(%r12)
	je	.L707
	cmpb	$0, 898(%r12)
	jne	.L707
	cmpb	$0, 896(%r12)
	je	.L930
	.p2align 4,,10
	.p2align 3
.L707:
	cmpw	$2, 698(%r12)
	leaq	EnvironmentOptions(%rip), %rdx
	movq	%r10, %rsi
	movq	%r12, %rdi
	jbe	.L710
	call	pqBuildStartupPacket3@PLT
	movq	%rax, %rbx
.L711:
	testq	%rbx, %rbx
	je	.L931
	movslq	80(%rsp), %rcx
	xorl	%esi, %esi
	movq	%rbx, %rdx
	movq	%r12, %rdi
	call	pqPacketSend@PLT
	testl	%eax, %eax
	je	.L713
	call	__errno_location@PLT
	movl	(%rax), %edi
	leaq	128(%rsp), %rsi
	movl	$256, %edx
	call	pg_strerror_r@PLT
	leaq	928(%r12), %rdi
	leaq	.LC64(%rip), %rsi
	movq	%rax, %rdx
	xorl	%eax, %eax
	call	appendPQExpBuffer@PLT
	movq	%rbx, %rdi
	call	free@PLT
	.p2align 4,,10
	.p2align 3
.L611:
	movl	$1, 336(%r12)
.L916:
	xorl	%eax, %eax
	jmp	.L602
	.p2align 4,,10
	.p2align 3
.L608:
	call	pqReadData@PLT
	testl	%eax, %eax
	js	.L611
	je	.L918
.L607:
	leaq	48(%rsp), %rax
	leaq	.L645(%rip), %rbx
	xorl	%r15d, %r15d
	xorl	%r13d, %r13d
	movq	%rax, (%rsp)
.L609:
	cmpb	$0, 720(%r12)
	jne	.L693
	cmpb	$0, 721(%r12)
	jne	.L613
	testb	%r13b, %r13b
	jne	.L621
	testb	%r15b, %r15b
	je	.L615
	.p2align 4,,10
	.p2align 3
.L617:
	movl	$1, %esi
	movq	%r12, %rdi
	call	pqDropConnection@PLT
	movq	%r12, %rdi
	call	pqDropServerData
	movl	$0, 340(%r12)
	movl	$0, 344(%r12)
	movq	%r12, %rdi
	call	pqClearAsyncResult@PLT
	movl	$8, 336(%r12)
.L642:
	movq	736(%r12), %r13
	testq	%r13, %r13
	je	.L932
	movl	16(%r13), %edx
	movq	24(%r13), %rsi
	leaq	560(%r12), %rdi
	call	memcpy@PLT
	movq	408(%r12), %rdi
	movl	16(%r13), %eax
	testq	%rdi, %rdi
	movl	%eax, 688(%r12)
	je	.L657
	call	free@PLT
	movq	$0, 408(%r12)
.L657:
	leaq	640(%rsp), %r14
	movq	%r12, %rdi
	movq	%r14, %rsi
	call	getHostaddr.constprop.14
	cmpb	$0, 640(%rsp)
	jne	.L933
.L658:
	movl	4(%r13), %edi
	xorl	%edx, %edx
	movl	$1, %esi
	call	socket@PLT
	cmpl	$-1, %eax
	movl	%eax, 416(%r12)
	je	.L934
	cmpl	$1, 4(%r13)
	je	.L663
	leaq	80(%rsp), %r10
	movl	$4, %r8d
	movl	$1, %edx
	movl	$6, %esi
	movl	%eax, %edi
	movl	$1, 80(%rsp)
	movq	%r10, %rcx
	call	setsockopt@PLT
	testl	%eax, %eax
	js	.L935
	movl	416(%r12), %eax
.L663:
	movl	%eax, %edi
	call	pg_set_noblock@PLT
	testb	%al, %al
	jne	.L665
	call	__errno_location@PLT
	movl	(%rax), %edi
	leaq	128(%rsp), %rsi
	movl	$256, %edx
	call	pg_strerror_r@PLT
	leaq	928(%r12), %rdi
	leaq	.LC43(%rip), %rsi
	movq	%rax, %rdx
	xorl	%eax, %eax
	call	appendPQExpBuffer@PLT
	movb	$1, 720(%r12)
	.p2align 4,,10
	.p2align 3
.L911:
	xorl	%r13d, %r13d
	xorl	%r15d, %r15d
.L693:
	movq	736(%r12), %rax
	testq	%rax, %rax
	je	.L620
	movq	40(%rax), %rax
	testq	%rax, %rax
	je	.L620
	movzbl	721(%r12), %r13d
	movq	%rax, 736(%r12)
	movb	$0, 720(%r12)
	testb	%r13b, %r13b
	jne	.L613
.L621:
	movq	168(%r12), %rax
	movl	$196608, 696(%r12)
	movb	$1, 760(%r12)
	cmpb	$100, (%rax)
	setne	897(%r12)
	cmpb	$97, (%rax)
	sete	898(%r12)
	jmp	.L617
	.p2align 4,,10
	.p2align 3
.L793:
	movl	$3, %eax
	jmp	.L602
.L930:
	xorl	%esi, %esi
	movl	$4, %ecx
	movq	%r10, %rdx
	movq	%r12, %rdi
	movl	$790024708, 80(%rsp)
	call	pqPacketSend@PLT
	testl	%eax, %eax
	jne	.L936
	movl	$7, 336(%r12)
	.p2align 4,,10
	.p2align 3
.L918:
	movl	$1, %eax
	jmp	.L602
	.p2align 4,,10
	.p2align 3
.L932:
	movzbl	720(%r12), %eax
	movb	$1, 721(%r12)
.L656:
	testb	%al, %al
	jne	.L911
	cmpb	$0, 721(%r12)
	je	.L615
.L614:
	xorl	%r13d, %r13d
	xorl	%r15d, %r15d
	.p2align 4,,10
	.p2align 3
.L613:
	movl	396(%r12), %eax
	addl	$1, %eax
	cmpl	392(%r12), %eax
	jge	.L611
	movq	728(%r12), %rsi
	movl	%eax, 396(%r12)
	testq	%rsi, %rsi
	je	.L624
	movl	744(%r12), %edi
	call	pg_freeaddrinfo_all@PLT
	movl	396(%r12), %eax
	movq	$0, 728(%r12)
	movq	$0, 736(%r12)
.L624:
	cltq
	movq	$0, 80(%rsp)
	movq	$1, 88(%rsp)
	leaq	(%rax,%rax,4), %rdx
	movq	400(%r12), %rax
	movq	$0, 96(%rsp)
	movq	$0, 104(%rsp)
	movq	$0, 112(%rsp)
	movq	$0, 120(%rsp)
	movl	$0, 744(%r12)
	leaq	(%rax,%rdx,8), %rbp
	movq	24(%rbp), %rdi
	testq	%rdi, %rdi
	je	.L625
	cmpb	$0, (%rdi)
	jne	.L626
.L625:
	movl	$5432, 40(%rsp)
	movl	$5432, %ecx
.L627:
	leaq	640(%rsp), %r14
	leaq	.LC35(%rip), %rdx
	xorl	%eax, %eax
	movl	$1024, %esi
	movq	%r14, %rdi
	call	pg_snprintf@PLT
	movl	0(%rbp), %eax
	cmpl	$1, %eax
	je	.L630
	jb	.L631
	cmpl	$2, %eax
	jne	.L937
	movl	$1, 84(%rsp)
	movl	$1, 744(%r12)
	leaq	.LC38(%rip), %rdx
	movq	8(%rbp), %rcx
	movl	40(%rsp), %r8d
	movl	$1024, %esi
	movq	%r14, %rdi
	xorl	%eax, %eax
	call	pg_snprintf@PLT
	movq	%r14, %rdx
.L636:
	movl	(%rdx), %ecx
	addq	$4, %rdx
	leal	-16843009(%rcx), %eax
	notl	%ecx
	andl	%ecx, %eax
	andl	$-2139062144, %eax
	je	.L636
	movl	%eax, %ecx
	shrl	$16, %ecx
	testl	$32896, %eax
	cmove	%ecx, %eax
	leaq	2(%rdx), %rcx
	cmove	%rcx, %rdx
	movl	%eax, %ecx
	addb	%al, %cl
	sbbq	$3, %rdx
	subq	%r14, %rdx
	cmpq	$107, %rdx
	jbe	.L638
	leaq	928(%r12), %rdi
	leaq	.LC39(%rip), %rsi
	movl	$107, %ecx
	movq	%r14, %rdx
	xorl	%eax, %eax
	call	appendPQExpBuffer@PLT
	jmp	.L609
	.p2align 4,,10
	.p2align 3
.L604:
	leaq	928(%rdi), %rdi
	leaq	.LC32(%rip), %rsi
	call	appendPQExpBufferStr@PLT
	movl	$1, 336(%r12)
	jmp	.L916
	.p2align 4,,10
	.p2align 3
.L620:
	movb	$1, 721(%r12)
	movb	$0, 720(%r12)
	jmp	.L613
	.p2align 4,,10
	.p2align 3
.L615:
	leaq	36(%rsp), %r14
	leaq	80(%rsp), %r15
.L641:
	movl	336(%r12), %edx
	cmpl	$12, %edx
	ja	.L643
	movl	%edx, %eax
	movslq	(%rbx,%rax,4), %rax
	addq	%rbx, %rax
	jmp	*%rax
	.section	.rodata
	.align 4
	.align 4
.L645:
	.long	.L643-.L645
	.long	.L643-.L645
	.long	.L644-.L645
	.long	.L646-.L645
	.long	.L647-.L645
	.long	.L648-.L645
	.long	.L649-.L645
	.long	.L650-.L645
	.long	.L642-.L645
	.long	.L651-.L645
	.long	.L652-.L645
	.long	.L611-.L645
	.long	.L654-.L645
	.text
	.p2align 4,,10
	.p2align 3
.L937:
	movq	728(%r12), %rdx
.L634:
	movq	%rdx, 736(%r12)
	movb	$0, 721(%r12)
	jmp	.L621
.L651:
	movq	(%rsp), %rbp
	leaq	928(%r12), %r13
	movq	%rbp, %rdi
	movq	%rbp, 8(%rsp)
	call	initPQExpBuffer@PLT
	movq	928(%r12), %rsi
	movq	%rbp, %rdi
	call	appendPQExpBufferStr@PLT
	cmpq	$0, 64(%rsp)
	jne	.L777
	leaq	.LC3(%rip), %rsi
	movq	%r13, %rdi
	xorl	%eax, %eax
	call	printfPQExpBuffer@PLT
	jmp	.L611
.L650:
	cmpb	$0, 896(%r12)
	jne	.L714
	movq	%r12, %rdi
	call	pqReadData@PLT
	testl	%eax, %eax
	js	.L611
	je	.L918
	movq	%r12, %rsi
	movq	%r15, %rdi
	call	pqGetc@PLT
	testl	%eax, %eax
	js	.L918
	movzbl	80(%rsp), %eax
	cmpb	$83, %al
	je	.L938
	cmpb	$78, %al
	je	.L939
	cmpb	$69, %al
	jne	.L723
	cmpb	$0, 720(%r12)
	movl	$4, 336(%r12)
	jne	.L911
	.p2align 4,,10
	.p2align 3
.L899:
	cmpb	$0, 721(%r12)
	je	.L641
	jmp	.L614
.L654:
	cmpl	$70399, 700(%r12)
	jle	.L764
	movq	272(%r12), %rsi
	testq	%rsi, %rsi
	je	.L764
	leaq	.LC72(%rip), %rdi
	movl	$11, %ecx
	repz cmpsb
	seta	%al
	sbbb	$0, %al
	testb	%al, %al
	jne	.L764
	leaq	48(%rsp), %rax
	movq	%rax, %rdi
	movq	%rax, %rbx
	movq	%rax, 8(%rsp)
	call	initPQExpBuffer@PLT
	movq	928(%r12), %rsi
	movq	%rbx, %rdi
	call	appendPQExpBufferStr@PLT
	cmpq	$0, 64(%rsp)
	leaq	928(%r12), %rdi
	jne	.L765
	leaq	.LC3(%rip), %rsi
	xorl	%eax, %eax
	call	printfPQExpBuffer@PLT
	jmp	.L611
.L652:
	movl	$0, 336(%r12)
	movq	%r12, %rdi
	call	PQconsumeInput@PLT
	testl	%eax, %eax
	je	.L611
	movq	%r12, %rdi
	call	PQisBusy@PLT
	testl	%eax, %eax
	jne	.L940
	movq	%r12, %rdi
	call	PQgetResult@PLT
	testq	%rax, %rax
	je	.L764
	movq	%rax, %rdi
	call	PQclear@PLT
	cmpb	$0, 720(%r12)
	movl	$10, 336(%r12)
	jne	.L911
	jmp	.L899
.L647:
	movl	820(%r12), %eax
	movq	%r12, %rsi
	movq	%r14, %rdi
	movl	%eax, 824(%r12)
	call	pqGetc@PLT
	testl	%eax, %eax
	jne	.L918
	movzbl	36(%rsp), %eax
	cmpb	$82, %al
	je	.L729
	cmpb	$69, %al
	jne	.L941
.L729:
	cmpw	$2, 698(%r12)
	ja	.L942
	cmpb	$82, %al
	movl	$8, 40(%rsp)
	movl	$8, %edx
	je	.L943
.L734:
	cmpb	$69, %al
	jne	.L915
	leal	-8(%rdx), %eax
	cmpl	$29992, %eax
	ja	.L944
	movslq	824(%r12), %rcx
	leal	-4(%rdx), %eax
	movl	828(%r12), %edx
	movl	%eax, 40(%rsp)
	subl	%ecx, %edx
	cmpl	%edx, %eax
	jg	.L789
	cmpw	$2, 698(%r12)
	jbe	.L742
	movl	$1, %esi
	movq	%r12, %rdi
	call	pqGetErrorNotice3@PLT
	testl	%eax, %eax
	jne	.L918
.L745:
	cmpb	$0, 705(%r12)
	movl	824(%r12), %eax
	movl	%eax, 820(%r12)
	je	.L746
	movslq	396(%r12), %rax
	movq	400(%r12), %rdx
	leaq	(%rax,%rax,4), %rax
	leaq	(%rdx,%rax,8), %rax
	cmpq	$0, 32(%rax)
	je	.L746
	movq	872(%r12), %rdi
	testq	%rdi, %rdi
	je	.L746
	movl	$67, %esi
	call	PQresultErrorField@PLT
	testq	%rax, %rax
	movq	%rax, %rsi
	je	.L746
	leaq	.LC68(%rip), %rdi
	movl	$6, %ecx
	repz cmpsb
	seta	%al
	sbbb	$0, %al
	testb	%al, %al
	jne	.L746
	movq	120(%r12), %rdx
	leaq	928(%r12), %rdi
	leaq	.LC69(%rip), %rsi
	call	appendPQExpBuffer@PLT
	.p2align 4,,10
	.p2align 3
.L746:
	movq	168(%r12), %rax
	movzbl	(%rax), %eax
	cmpb	$97, %al
	je	.L945
	cmpb	$112, %al
	jne	.L611
	cmpb	$0, 896(%r12)
	je	.L611
.L927:
	cmpb	$0, 897(%r12)
	je	.L611
	cmpb	$0, 898(%r12)
	jne	.L611
	cmpb	$0, 720(%r12)
	movb	$0, 897(%r12)
	je	.L727
.L739:
	xorl	%r13d, %r13d
	movl	$1, %r15d
	jmp	.L693
.L648:
	movq	%r12, %rdi
	call	PQisBusy@PLT
	testl	%eax, %eax
	jne	.L918
	movq	%r12, %rdi
	call	PQgetResult@PLT
	testq	%rax, %rax
	movq	%rax, %r13
	je	.L754
	cmpl	$7, 40(%rax)
	jne	.L946
	cmpb	$0, 760(%r12)
	je	.L756
	cmpq	$0, 72(%r12)
	je	.L947
.L758:
	movl	$67, %esi
	movq	%r13, %rdi
	call	PQresultErrorField@PLT
	testq	%rax, %rax
	movq	%rax, %rsi
	je	.L756
	leaq	.LC71(%rip), %rdi
	movl	$6, %ecx
	repz cmpsb
	seta	%al
	sbbb	$0, %al
	testb	%al, %al
	jne	.L756
	movq	%r13, %rdi
	call	PQclear@PLT
	cmpb	$0, 720(%r12)
	movb	$0, 760(%r12)
	jne	.L739
.L727:
	xorl	%r13d, %r13d
	cmpb	$0, 721(%r12)
	movl	$1, %r15d
	je	.L617
	jmp	.L613
.L644:
	movl	416(%r12), %edi
	leaq	28(%rsp), %rcx
	leaq	80(%rsp), %r8
	movl	$4, %edx
	movl	$1, %esi
	movl	$4, 80(%rsp)
	call	getsockopt@PLT
	cmpl	$-1, %eax
	je	.L948
	movl	28(%rsp), %esi
	testl	%esi, %esi
	je	.L696
.L913:
	movq	%r12, %rdi
	call	connectFailureMessage
.L912:
	movb	$1, 720(%r12)
	jmp	.L911
.L643:
	leaq	928(%r12), %rdi
	leaq	.LC77(%rip), %rsi
	xorl	%eax, %eax
	call	appendPQExpBuffer@PLT
	jmp	.L611
	.p2align 4,,10
	.p2align 3
.L626:
	leaq	40(%rsp), %rsi
	leaq	.LC33(%rip), %rcx
	movq	%r12, %rdx
	call	parse_int_param
	testb	%al, %al
	je	.L611
	movl	40(%rsp), %ecx
	leal	-1(%rcx), %eax
	cmpl	$65534, %eax
	jbe	.L627
	movq	24(%rbp), %rdx
	leaq	928(%r12), %rdi
	leaq	.LC34(%rip), %rsi
	xorl	%eax, %eax
	call	appendPQExpBuffer@PLT
	jmp	.L609
	.p2align 4,,10
	.p2align 3
.L631:
	leaq	80(%rsp), %r10
	movq	8(%rbp), %rdi
	leaq	728(%r12), %rcx
	movq	%r14, %rsi
	movq	%r10, %rdx
	call	pg_getaddrinfo_all@PLT
	testl	%eax, %eax
	jne	.L633
	movq	728(%r12), %rdx
	testq	%rdx, %rdx
	jne	.L634
.L633:
	movl	%eax, %edi
	call	gai_strerror@PLT
	movq	8(%rbp), %rdx
	leaq	928(%r12), %rdi
	leaq	.LC36(%rip), %rsi
	movq	%rax, %rcx
	xorl	%eax, %eax
	call	appendPQExpBuffer@PLT
	jmp	.L609
	.p2align 4,,10
	.p2align 3
.L630:
	movl	$4, 80(%rsp)
	leaq	80(%rsp), %r10
	movq	16(%rbp), %rdi
	leaq	728(%r12), %rcx
	movq	%r14, %rsi
	movq	%r10, %rdx
	call	pg_getaddrinfo_all@PLT
	testl	%eax, %eax
	jne	.L635
	movq	728(%r12), %rdx
	testq	%rdx, %rdx
	jne	.L634
.L635:
	movl	%eax, %edi
	call	gai_strerror@PLT
	movq	16(%rbp), %rdx
	leaq	928(%r12), %rdi
	leaq	.LC37(%rip), %rsi
	movq	%rax, %rcx
	xorl	%eax, %eax
	call	appendPQExpBuffer@PLT
	jmp	.L609
	.p2align 4,,10
	.p2align 3
.L777:
	movq	%r13, %rdi
	call	resetPQExpBuffer@PLT
	movl	$0, 336(%r12)
	movq	%r12, %rdi
	call	PQconsumeInput@PLT
	testl	%eax, %eax
	je	.L778
	movq	%r12, %rdi
	call	PQisBusy@PLT
	testl	%eax, %eax
	jne	.L766
	movq	%r12, %rdi
	call	PQgetResult@PLT
	testq	%rax, %rax
	je	.L779
	movq	%rax, %rdi
	movq	%rax, 8(%rsp)
	call	PQresultStatus@PLT
	cmpl	$2, %eax
	movq	8(%rsp), %rcx
	je	.L949
.L780:
	movq	%rcx, %rdi
	call	PQclear@PLT
.L779:
	movq	(%rsp), %rsi
	movq	%r12, %rdi
	call	restoreErrorMessage
	movslq	396(%r12), %rax
	leaq	(%rax,%rax,4), %rdx
	movq	400(%r12), %rax
	leaq	(%rax,%rdx,8), %rax
	movq	24(%rax), %rcx
	cmpl	$1, (%rax)
	movq	16(%rax), %rdx
	cmovne	8(%rax), %rdx
	testq	%rcx, %rcx
	je	.L798
	cmpb	$0, (%rcx)
	leaq	.LC26(%rip), %rax
	cmove	%rax, %rcx
.L787:
	leaq	.LC76(%rip), %rsi
	xorl	%eax, %eax
	movq	%r13, %rdi
	call	appendPQExpBuffer@PLT
	cmpl	$-1, 416(%r12)
	movl	$0, 336(%r12)
	je	.L912
	movq	%r12, %rdi
	call	sendTerminateConn.part.7
	jmp	.L912
	.p2align 4,,10
	.p2align 3
.L764:
	movq	728(%r12), %rsi
	testq	%rsi, %rsi
	je	.L768
	movl	744(%r12), %edi
	call	pg_freeaddrinfo_all@PLT
	movq	$0, 728(%r12)
	movq	$0, 736(%r12)
.L768:
	movl	$0, 336(%r12)
	movl	$3, %eax
	jmp	.L602
.L700:
	movl	32(%rsp), %edi
	leaq	640(%rsp), %rdx
	leaq	40(%rsp), %r8
	movq	%r10, %rsi
	movl	$8192, %ecx
	movq	%r10, (%rsp)
	call	pqGetpwuid@PLT
	movq	40(%rsp), %rdx
	movq	(%rsp), %r10
	testq	%rdx, %rdx
	je	.L950
	movq	(%rdx), %rbp
	movq	224(%r12), %rbx
	movq	%r10, (%rsp)
	movq	%rbx, %rsi
	movq	%rbp, %rdi
	call	strcmp@PLT
	testl	%eax, %eax
	movq	(%rsp), %r10
	jne	.L951
	movzwl	560(%r12), %edx
	.p2align 4,,10
	.p2align 3
.L698:
	cmpw	$1, %dx
	jne	.L699
	movb	$0, 897(%r12)
	jmp	.L707
.L938:
	movl	824(%r12), %eax
	movq	%r12, %rdi
	movl	%eax, 820(%r12)
	call	pqsecure_initialize@PLT
	testl	%eax, %eax
	jne	.L611
	.p2align 4,,10
	.p2align 3
.L714:
	movq	%r12, %rdi
	call	pqsecure_open_client@PLT
	cmpl	$3, %eax
	je	.L917
	testl	%eax, %eax
	jne	.L602
	movq	168(%r12), %rax
	cmpb	$112, (%rax)
	je	.L927
	jmp	.L611
	.p2align 4,,10
	.p2align 3
.L638:
	leaq	80(%rsp), %r10
	leaq	728(%r12), %rcx
	xorl	%edi, %edi
	movq	%r14, %rsi
	movq	%r10, %rdx
	call	pg_getaddrinfo_all@PLT
	testl	%eax, %eax
	jne	.L639
	movq	728(%r12), %rdx
	testq	%rdx, %rdx
	jne	.L634
.L639:
	movl	%eax, %edi
	call	gai_strerror@PLT
	leaq	928(%r12), %rdi
	leaq	.LC40(%rip), %rsi
	movq	%rax, %rcx
	movq	%r14, %rdx
	xorl	%eax, %eax
	call	appendPQExpBuffer@PLT
	jmp	.L609
	.p2align 4,,10
	.p2align 3
.L770:
	cmpb	$0, 720(%r12)
	movl	$12, 336(%r12)
	jne	.L911
	jmp	.L899
	.p2align 4,,10
	.p2align 3
.L710:
	call	pqBuildStartupPacket2@PLT
	movq	%rax, %rbx
	jmp	.L711
	.p2align 4,,10
	.p2align 3
.L915:
	leal	-4(%rdx), %eax
.L733:
	movslq	824(%r12), %rcx
	movl	828(%r12), %edx
	movl	%eax, 40(%rsp)
	subl	%ecx, %edx
	cmpl	%edx, %eax
	jg	.L789
	movb	$1, 704(%r12)
	movq	%r12, %rdx
	movl	$4, %esi
	movq	%r15, %rdi
	call	pqGetInt@PLT
	testl	%eax, %eax
	jne	.L918
	movl	40(%rsp), %eax
	movl	80(%rsp), %edi
	leal	-4(%rax), %esi
	cmpl	$5, %edi
	movl	%esi, 40(%rsp)
	jne	.L750
	cmpw	$2, 698(%r12)
	ja	.L750
	movslq	824(%r12), %rcx
	movl	828(%r12), %edx
	movl	%eax, 40(%rsp)
	subl	%ecx, %edx
	cmpl	$3, %edx
	jle	.L952
	movl	%eax, %esi
.L750:
	movq	%r12, %rdx
	call	pg_fe_sendauth@PLT
	movq	928(%r12), %rdi
	movl	%eax, %r13d
	call	strlen@PLT
	movq	%rax, 936(%r12)
	movl	824(%r12), %eax
	testl	%r13d, %r13d
	movl	%eax, 820(%r12)
	jne	.L611
	movq	%r12, %rdi
	call	pqFlush@PLT
	testl	%eax, %eax
	jne	.L611
	movl	80(%rsp), %eax
	testl	%eax, %eax
	jne	.L921
	movl	$5, 336(%r12)
	movl	$1, 340(%r12)
.L921:
	cmpb	$0, 720(%r12)
	jne	.L911
	jmp	.L899
	.p2align 4,,10
	.p2align 3
.L942:
	leaq	40(%rsp), %rdi
	movq	%r12, %rdx
	movl	$4, %esi
	call	pqGetInt@PLT
	testl	%eax, %eax
	jne	.L918
	movzbl	36(%rsp), %eax
	movl	40(%rsp), %edx
	cmpb	$82, %al
	jne	.L734
	leal	-8(%rdx), %eax
	cmpl	$1992, %eax
	jbe	.L915
	leaq	928(%r12), %rdi
	leaq	.LC67(%rip), %rsi
	movl	$82, %edx
	xorl	%eax, %eax
	call	appendPQExpBuffer@PLT
	jmp	.L611
	.p2align 4,,10
	.p2align 3
.L947:
	cmpq	$0, 80(%r12)
	jne	.L758
.L756:
	movq	936(%r12), %rax
	testq	%rax, %rax
	je	.L760
	movq	928(%r12), %rdx
	cmpb	$10, -1(%rdx,%rax)
	je	.L761
.L760:
	leaq	928(%r12), %rdi
	movl	$10, %esi
	call	appendPQExpBufferChar@PLT
.L761:
	movq	%r13, %rdi
	call	PQclear@PLT
	jmp	.L611
	.p2align 4,,10
	.p2align 3
.L713:
	movq	%rbx, %rdi
	call	free@PLT
	movl	$4, 336(%r12)
	jmp	.L918
	.p2align 4,,10
	.p2align 3
.L754:
	cmpw	$2, 698(%r12)
	ja	.L770
	leaq	EnvironmentOptions(%rip), %rax
	movl	$6, 336(%r12)
	movl	$0, 748(%r12)
	movq	%rax, 752(%r12)
	movl	$2, %eax
	jmp	.L602
	.p2align 4,,10
	.p2align 3
.L798:
	leaq	.LC26(%rip), %rcx
	jmp	.L787
.L931:
	leaq	928(%r12), %rdi
	leaq	.LC3(%rip), %rsi
	xorl	%eax, %eax
	call	printfPQExpBuffer@PLT
	jmp	.L611
.L769:
	movl	$6, 336(%r12)
	movl	$2, %eax
	jmp	.L602
.L943:
	movl	$4, %eax
	jmp	.L733
.L742:
	leaq	928(%r12), %rdi
	movq	%r12, %rsi
	call	pqGets_append@PLT
	testl	%eax, %eax
	je	.L745
	jmp	.L918
	.p2align 4,,10
	.p2align 3
.L934:
	cmpq	$0, 40(%r13)
	jne	.L912
	movl	396(%r12), %eax
	addl	$1, %eax
	cmpl	392(%r12), %eax
	jl	.L912
	call	__errno_location@PLT
	movl	(%rax), %edi
	leaq	128(%rsp), %rsi
	movl	$256, %edx
	call	pg_strerror_r@PLT
	leaq	928(%r12), %rdi
	leaq	.LC41(%rip), %rsi
	movq	%rax, %rdx
	xorl	%eax, %eax
	call	appendPQExpBuffer@PLT
	jmp	.L611
	.p2align 4,,10
	.p2align 3
.L933:
	movq	%r14, %rdi
	call	strdup@PLT
	movq	%rax, 408(%r12)
	jmp	.L658
.L945:
	cmpb	$0, 896(%r12)
	jne	.L611
	cmpb	$0, 897(%r12)
	je	.L611
	cmpb	$0, 898(%r12)
	je	.L611
	cmpb	$0, 720(%r12)
	movb	$0, 898(%r12)
	jne	.L739
	jmp	.L727
	.p2align 4,,10
	.p2align 3
.L665:
	movl	416(%r12), %edi
	xorl	%eax, %eax
	movl	$1, %edx
	movl	$2, %esi
	call	fcntl@PLT
	cmpl	$-1, %eax
	je	.L953
	cmpl	$1, 4(%r13)
	je	.L667
	movq	136(%r12), %rdi
	movl	$1, 40(%rsp)
	testq	%rdi, %rdi
	je	.L669
	leaq	80(%rsp), %r10
	movl	$10, %edx
	movq	%r10, %rsi
	call	strtol@PLT
	movq	80(%rsp), %rdx
	cmpb	$0, (%rdx)
	jne	.L954
	testl	%eax, %eax
	jne	.L669
.L667:
	movb	$0, 706(%r12)
	movb	$1, 707(%r12)
	movl	416(%r12), %edi
	movl	16(%r13), %edx
	movq	24(%r13), %rsi
	call	connect@PLT
	testl	%eax, %eax
	js	.L955
	movl	$2, 336(%r12)
	movzbl	720(%r12), %eax
	jmp	.L656
.L935:
	call	__errno_location@PLT
	movl	(%rax), %edi
	leaq	384(%rsp), %rsi
	movl	$256, %edx
	call	pg_strerror_r@PLT
	leaq	928(%r12), %rdi
	leaq	.LC42(%rip), %rsi
	movq	%rax, %rdx
	xorl	%eax, %eax
	call	appendPQExpBuffer@PLT
	movb	$1, 720(%r12)
	jmp	.L911
.L765:
	call	resetPQExpBuffer@PLT
	leaq	.LC73(%rip), %rsi
	movl	$0, 336(%r12)
	movq	%r12, %rdi
	call	PQsendQuery@PLT
	testl	%eax, %eax
	jne	.L766
.L778:
	movq	8(%rsp), %rsi
	movq	%r12, %rdi
	call	restoreErrorMessage
	jmp	.L611
.L944:
	movl	820(%r12), %eax
	leaq	928(%r12), %r13
	movq	%r12, %rsi
	movq	%r13, %rdi
	addl	$1, %eax
	movl	%eax, 824(%r12)
	call	pqGets_append@PLT
	testl	%eax, %eax
	jne	.L918
	movl	824(%r12), %eax
	movl	$10, %esi
	movq	%r13, %rdi
	movl	%eax, 820(%r12)
	call	appendPQExpBufferChar@PLT
	cmpw	$2, 698(%r12)
	jbe	.L611
	cmpb	$0, 720(%r12)
	movl	$131072, 696(%r12)
	jne	.L739
	jmp	.L727
.L949:
	movq	%rcx, %rdi
	call	PQntuples@PLT
	cmpl	$1, %eax
	movq	8(%rsp), %rcx
	jne	.L780
	movq	%rcx, %rdi
	xorl	%edx, %edx
	xorl	%esi, %esi
	call	PQgetvalue@PLT
	leaq	.LC74(%rip), %rsi
	movq	%rax, %rdi
	movl	$2, %edx
	call	strncmp@PLT
	movq	8(%rsp), %rcx
	testl	%eax, %eax
	movq	%rcx, %rdi
	jne	.L781
	call	PQclear@PLT
	movq	(%rsp), %rsi
	movq	%r12, %rdi
	call	restoreErrorMessage
	movslq	396(%r12), %rax
	imulq	$40, %rax, %rax
	addq	400(%r12), %rax
	movq	24(%rax), %rcx
	cmpl	$1, (%rax)
	movq	16(%rax), %rdx
	cmovne	8(%rax), %rdx
	testq	%rcx, %rcx
	je	.L797
	cmpb	$0, (%rcx)
	jne	.L783
.L797:
	leaq	.LC26(%rip), %rcx
.L783:
	leaq	.LC75(%rip), %rsi
	xorl	%eax, %eax
	movq	%r13, %rdi
	call	appendPQExpBuffer@PLT
	cmpl	$-1, 416(%r12)
	movl	$0, 336(%r12)
	je	.L784
	movq	%r12, %rdi
	call	sendTerminateConn.part.7
.L784:
	cmpb	$0, 720(%r12)
	movb	$1, 721(%r12)
	jne	.L911
	jmp	.L614
	.p2align 4,,10
	.p2align 3
.L929:
	movq	%r10, (%rsp)
	call	__errno_location@PLT
	movl	$0, (%rax)
	movl	416(%r12), %edi
	leaq	36(%rsp), %rdx
	leaq	32(%rsp), %rsi
	movq	%rax, %rbx
	call	getpeereid@PLT
	testl	%eax, %eax
	movq	(%rsp), %r10
	je	.L700
	movl	(%rbx), %edi
	leaq	928(%r12), %rbx
	cmpl	$38, %edi
	je	.L956
	leaq	128(%rsp), %rsi
	movl	$256, %edx
	call	pg_strerror_r@PLT
	leaq	.LC59(%rip), %rsi
	movq	%rax, %rdx
	movq	%rbx, %rdi
	xorl	%eax, %eax
	call	appendPQExpBuffer@PLT
	jmp	.L611
.L766:
	movq	8(%rsp), %rsi
	movl	$9, 336(%r12)
	movq	%r12, %rdi
	call	restoreErrorMessage
	movl	$1, %eax
	jmp	.L602
.L939:
	movl	824(%r12), %eax
	movl	%eax, 820(%r12)
	movq	168(%r12), %rax
	movzbl	(%rax), %eax
	andl	$-5, %eax
	cmpb	$114, %al
	je	.L957
	movb	$0, 897(%r12)
.L917:
	movl	$3, 336(%r12)
	movl	$2, %eax
	jmp	.L602
.L936:
	call	__errno_location@PLT
	movl	(%rax), %edi
	leaq	128(%rsp), %rsi
	movl	$256, %edx
	call	pg_strerror_r@PLT
	leaq	928(%r12), %rdi
	leaq	.LC63(%rip), %rsi
	movq	%rax, %rdx
	xorl	%eax, %eax
	call	appendPQExpBuffer@PLT
	jmp	.L611
.L955:
	call	__errno_location@PLT
	movl	(%rax), %esi
	cmpl	$115, %esi
	je	.L800
	cmpl	$4, %esi
	jne	.L913
.L800:
	movl	$2, 336(%r12)
	movl	$2, %eax
	jmp	.L602
.L953:
	call	__errno_location@PLT
	movl	(%rax), %edi
	leaq	128(%rsp), %rsi
	movl	$256, %edx
	call	pg_strerror_r@PLT
	leaq	928(%r12), %rdi
	leaq	.LC44(%rip), %rsi
	movq	%rax, %rdx
	xorl	%eax, %eax
	call	appendPQExpBuffer@PLT
	movb	$1, 720(%r12)
	jmp	.L911
.L941:
	leaq	928(%r12), %rdi
	leaq	.LC67(%rip), %rsi
	movsbl	%al, %edx
	xorl	%eax, %eax
	call	appendPQExpBuffer@PLT
	jmp	.L611
.L940:
	movl	$10, 336(%r12)
	jmp	.L918
.L696:
	movl	416(%r12), %edi
	leaq	424(%r12), %rsi
	leaq	552(%r12), %rdx
	movl	$128, 552(%r12)
	call	getsockname@PLT
	testl	%eax, %eax
	jns	.L917
	call	__errno_location@PLT
	movl	(%rax), %edi
	leaq	128(%rsp), %rsi
	movl	$256, %edx
	call	pg_strerror_r@PLT
	leaq	928(%r12), %rdi
	leaq	.LC57(%rip), %rsi
	movq	%rax, %rdx
	xorl	%eax, %eax
	call	appendPQExpBuffer@PLT
	jmp	.L611
.L948:
	call	__errno_location@PLT
	movl	(%rax), %edi
	leaq	128(%rsp), %rsi
	movl	$256, %edx
	call	pg_strerror_r@PLT
	leaq	928(%r12), %rdi
	leaq	.LC56(%rip), %rsi
	movq	%rax, %rdx
	xorl	%eax, %eax
	call	appendPQExpBuffer@PLT
	jmp	.L611
.L956:
	leaq	.LC58(%rip), %rsi
	movq	%rbx, %rdi
	call	appendPQExpBufferStr@PLT
	jmp	.L611
.L789:
	cltq
	movq	%r12, %rsi
	leaq	(%rcx,%rax), %rdi
	call	pqCheckInBufferSpace@PLT
	testl	%eax, %eax
	jne	.L611
	jmp	.L918
.L669:
	movl	416(%r12), %edi
	leaq	40(%rsp), %rcx
	movl	$4, %r8d
	movl	$9, %edx
	movl	$1, %esi
	call	setsockopt@PLT
	testl	%eax, %eax
	js	.L958
	movq	144(%r12), %rdi
	testq	%rdi, %rdi
	je	.L674
	leaq	80(%rsp), %r10
	leaq	.LC48(%rip), %rcx
	movq	%r12, %rdx
	movq	%r10, %rsi
	movq	%r10, 8(%rsp)
	call	parse_int_param
	testb	%al, %al
	je	.L912
	cmpl	$0, 80(%rsp)
	movq	8(%rsp), %r10
	js	.L959
.L676:
	movl	416(%r12), %edi
	movl	$4, %r8d
	movq	%r10, %rcx
	movl	$4, %edx
	movl	$6, %esi
	call	setsockopt@PLT
	testl	%eax, %eax
	js	.L960
.L674:
	movq	152(%r12), %rdi
	testq	%rdi, %rdi
	je	.L677
	leaq	80(%rsp), %r10
	leaq	.LC50(%rip), %rcx
	movq	%r12, %rdx
	movq	%r10, %rsi
	movq	%r10, 8(%rsp)
	call	parse_int_param
	testb	%al, %al
	je	.L912
	cmpl	$0, 80(%rsp)
	movq	8(%rsp), %r10
	jns	.L680
	movl	$0, 80(%rsp)
.L680:
	movl	416(%r12), %edi
	movl	$4, %r8d
	movq	%r10, %rcx
	movl	$5, %edx
	movl	$6, %esi
	call	setsockopt@PLT
	testl	%eax, %eax
	js	.L961
.L677:
	movq	160(%r12), %rdi
	testq	%rdi, %rdi
	je	.L681
	leaq	80(%rsp), %r10
	leaq	.LC52(%rip), %rcx
	movq	%r12, %rdx
	movq	%r10, %rsi
	movq	%r10, 8(%rsp)
	call	parse_int_param
	testb	%al, %al
	je	.L912
	cmpl	$0, 80(%rsp)
	movq	8(%rsp), %r10
	jns	.L684
	movl	$0, 80(%rsp)
.L684:
	movl	416(%r12), %edi
	movl	$4, %r8d
	movq	%r10, %rcx
	movl	$6, %edx
	movl	$6, %esi
	call	setsockopt@PLT
	testl	%eax, %eax
	js	.L962
.L681:
	movq	48(%r12), %rdi
	testq	%rdi, %rdi
	je	.L667
	leaq	80(%rsp), %r10
	leaq	.LC54(%rip), %rcx
	movq	%r12, %rdx
	movq	%r10, %rsi
	movq	%r10, 8(%rsp)
	call	parse_int_param
	testb	%al, %al
	je	.L912
	cmpl	$0, 80(%rsp)
	movq	8(%rsp), %r10
	jns	.L688
	movl	$0, 80(%rsp)
.L688:
	movl	416(%r12), %edi
	movl	$4, %r8d
	movq	%r10, %rcx
	movl	$18, %edx
	movl	$6, %esi
	call	setsockopt@PLT
	testl	%eax, %eax
	jns	.L667
	call	__errno_location@PLT
	movl	(%rax), %edi
	leaq	384(%rsp), %rsi
	movl	$256, %edx
	call	pg_strerror_r@PLT
	leaq	928(%r12), %rdi
	leaq	.LC55(%rip), %rdx
	leaq	.LC47(%rip), %rsi
	movq	%rax, %rcx
	xorl	%eax, %eax
	call	appendPQExpBuffer@PLT
	jmp	.L912
.L954:
	leaq	928(%r12), %rdi
	leaq	.LC45(%rip), %rsi
	call	appendPQExpBufferStr@PLT
	jmp	.L912
.L946:
	leaq	928(%r12), %rdi
	leaq	.LC70(%rip), %rsi
	call	appendPQExpBufferStr@PLT
	jmp	.L756
.L781:
	call	PQclear@PLT
	movq	(%rsp), %rdi
	call	termPQExpBuffer@PLT
	movl	$10, 336(%r12)
	jmp	.L921
.L928:
	call	__stack_chk_fail@PLT
.L723:
	leaq	928(%r12), %rdi
	leaq	.LC66(%rip), %rsi
	movsbl	%al, %edx
	xorl	%eax, %eax
	call	appendPQExpBuffer@PLT
	jmp	.L611
.L958:
	call	__errno_location@PLT
	movl	(%rax), %edi
	leaq	128(%rsp), %rsi
	movl	$256, %edx
	call	pg_strerror_r@PLT
	leaq	928(%r12), %rdi
	leaq	.LC46(%rip), %rdx
	leaq	.LC47(%rip), %rsi
	movq	%rax, %rcx
	xorl	%eax, %eax
	call	appendPQExpBuffer@PLT
	jmp	.L912
.L951:
	leaq	928(%r12), %rdi
	leaq	.LC62(%rip), %rsi
	movq	%rbp, %rcx
	movq	%rbx, %rdx
	xorl	%eax, %eax
	call	appendPQExpBuffer@PLT
	jmp	.L611
.L950:
	testl	%eax, %eax
	leaq	928(%r12), %rbx
	je	.L704
	leaq	128(%rsp), %rsi
	movl	$256, %edx
	movl	%eax, %edi
	call	pg_strerror_r@PLT
	movl	32(%rsp), %edx
	leaq	.LC60(%rip), %rsi
	movq	%rax, %rcx
	movq	%rbx, %rdi
	xorl	%eax, %eax
	call	appendPQExpBuffer@PLT
	jmp	.L611
.L952:
	leaq	4(%rcx), %rdi
	movq	%r12, %rsi
	call	pqCheckInBufferSpace@PLT
	testl	%eax, %eax
	jne	.L611
	jmp	.L918
.L957:
	leaq	928(%r12), %rdi
	leaq	.LC65(%rip), %rsi
	call	appendPQExpBufferStr@PLT
	jmp	.L611
.L704:
	movl	32(%rsp), %edx
	leaq	.LC61(%rip), %rsi
	movq	%rbx, %rdi
	xorl	%eax, %eax
	call	appendPQExpBuffer@PLT
	jmp	.L611
.L961:
	call	__errno_location@PLT
	movl	(%rax), %edi
	leaq	384(%rsp), %rsi
	movl	$256, %edx
	call	pg_strerror_r@PLT
	leaq	928(%r12), %rdi
	leaq	.LC51(%rip), %rdx
	leaq	.LC47(%rip), %rsi
	movq	%rax, %rcx
	xorl	%eax, %eax
	call	appendPQExpBuffer@PLT
	jmp	.L912
.L962:
	call	__errno_location@PLT
	movl	(%rax), %edi
	leaq	384(%rsp), %rsi
	movl	$256, %edx
	call	pg_strerror_r@PLT
	leaq	928(%r12), %rdi
	leaq	.LC53(%rip), %rdx
	leaq	.LC47(%rip), %rsi
	movq	%rax, %rcx
	xorl	%eax, %eax
	call	appendPQExpBuffer@PLT
	jmp	.L912
.L960:
	call	__errno_location@PLT
	movl	(%rax), %edi
	leaq	384(%rsp), %rsi
	movl	$256, %edx
	call	pg_strerror_r@PLT
	leaq	928(%r12), %rdi
	leaq	.LC49(%rip), %rdx
	leaq	.LC47(%rip), %rsi
	movq	%rax, %rcx
	xorl	%eax, %eax
	call	appendPQExpBuffer@PLT
	jmp	.L912
.L959:
	movl	$0, 80(%rsp)
	jmp	.L676
	.cfi_endproc
.LFE937:
	.size	PQconnectPoll, .-PQconnectPoll
	.section	.rodata.str1.1
.LC78:
	.string	"connect_timeout"
	.text
	.p2align 4,,15
	.type	connectDBComplete, @function
connectDBComplete:
.LFB934:
	.cfi_startproc
	pushq	%r14
	.cfi_def_cfa_offset 16
	.cfi_offset 14, -16
	pushq	%r13
	.cfi_def_cfa_offset 24
	.cfi_offset 13, -24
	pushq	%r12
	.cfi_def_cfa_offset 32
	.cfi_offset 12, -32
	pushq	%rbp
	.cfi_def_cfa_offset 40
	.cfi_offset 6, -40
	pushq	%rbx
	.cfi_def_cfa_offset 48
	.cfi_offset 3, -48
	subq	$16, %rsp
	.cfi_def_cfa_offset 64
	movq	%fs:40, %rax
	movq	%rax, 8(%rsp)
	xorl	%eax, %eax
	testq	%rdi, %rdi
	movl	$0, 4(%rsp)
	je	.L963
	cmpl	$1, 336(%rdi)
	movq	%rdi, %rbx
	je	.L963
	movq	40(%rdi), %rdi
	testq	%rdi, %rdi
	je	.L982
	leaq	4(%rsp), %rsi
	leaq	.LC78(%rip), %rcx
	movq	%rbx, %rdx
	call	parse_int_param
	testb	%al, %al
	je	.L990
	movl	4(%rsp), %eax
	testl	%eax, %eax
	jle	.L967
	cmpl	$1, %eax
	jne	.L965
	movl	$2, 4(%rsp)
	movl	$2, %eax
	jmp	.L965
	.p2align 4,,10
	.p2align 3
.L982:
	xorl	%eax, %eax
.L965:
	xorl	%r14d, %r14d
	movl	$-2, %r13d
	movq	$-1, %r12
	movl	$2, %ebp
.L968:
	testl	%eax, %eax
	jle	.L969
	cmpl	%r13d, 396(%rbx)
	jne	.L970
	cmpq	%r14, 736(%rbx)
	je	.L969
.L970:
	xorl	%edi, %edi
	call	time@PLT
	movslq	4(%rsp), %r12
	movl	396(%rbx), %r13d
	movq	736(%rbx), %r14
	addq	%rax, %r12
.L969:
	cmpl	$1, %ebp
	je	.L972
	cmpl	$2, %ebp
	jne	.L990
	xorl	%edi, %edi
	movq	%r12, %rcx
	movq	%rbx, %rdx
	movl	$1, %esi
	call	pqWaitTimed@PLT
	cmpl	$-1, %eax
	je	.L990
.L975:
	cmpl	$1, %eax
	jne	.L976
	movb	$1, 720(%rbx)
	movl	$8, 336(%rbx)
.L976:
	movq	%rbx, %rdi
	call	PQconnectPoll@PLT
	cmpl	$3, %eax
	movl	%eax, %ebp
	jne	.L993
	leaq	928(%rbx), %rdi
	call	resetPQExpBuffer@PLT
	movl	$1, %eax
	.p2align 4,,10
	.p2align 3
.L963:
	movq	8(%rsp), %rdx
	xorq	%fs:40, %rdx
	jne	.L994
	addq	$16, %rsp
	.cfi_remember_state
	.cfi_def_cfa_offset 48
	popq	%rbx
	.cfi_def_cfa_offset 40
	popq	%rbp
	.cfi_def_cfa_offset 32
	popq	%r12
	.cfi_def_cfa_offset 24
	popq	%r13
	.cfi_def_cfa_offset 16
	popq	%r14
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L972:
	.cfi_restore_state
	xorl	%esi, %esi
	movq	%r12, %rcx
	movq	%rbx, %rdx
	movl	$1, %edi
	call	pqWaitTimed@PLT
	cmpl	$-1, %eax
	jne	.L975
	.p2align 4,,10
	.p2align 3
.L990:
	movl	$1, 336(%rbx)
	xorl	%eax, %eax
	jmp	.L963
	.p2align 4,,10
	.p2align 3
.L993:
	movl	4(%rsp), %eax
	jmp	.L968
	.p2align 4,,10
	.p2align 3
.L967:
	movl	$0, 4(%rsp)
	xorl	%eax, %eax
	jmp	.L965
.L994:
	call	__stack_chk_fail@PLT
	.cfi_endproc
.LFE934:
	.size	connectDBComplete, .-connectDBComplete
	.section	.rodata.str1.1
.LC79:
	.string	"57P03"
	.text
	.p2align 4,,15
	.type	internal_ping, @function
internal_ping:
.LFB938:
	.cfi_startproc
	testq	%rdi, %rdi
	movl	$3, %eax
	je	.L1001
	cmpb	$0, 366(%rdi)
	je	.L1006
	pushq	%rbx
	.cfi_def_cfa_offset 16
	.cfi_offset 3, -16
	cmpl	$1, 336(%rdi)
	movq	%rdi, %rbx
	jne	.L997
.L1000:
	cmpb	$0, 704(%rbx)
	jne	.L998
	addq	$360, %rbx
	movq	%rbx, %rdi
	call	strlen@PLT
	movq	%rax, %rdx
	movl	$2, %eax
	cmpq	$5, %rdx
	jne	.L995
	leaq	.LC79(%rip), %rdi
	movl	$6, %ecx
	movq	%rbx, %rsi
	repz cmpsb
	seta	%al
	sbbb	$0, %al
	testb	%al, %al
	sete	%al
	movzbl	%al, %eax
.L995:
	popq	%rbx
	.cfi_remember_state
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L997:
	.cfi_restore_state
	call	connectDBComplete
	cmpl	$1, 336(%rbx)
	je	.L1000
.L998:
	xorl	%eax, %eax
	popq	%rbx
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L1006:
	.cfi_restore 3
	rep ret
	.p2align 4,,10
	.p2align 3
.L1001:
	rep ret
	.cfi_endproc
.LFE938:
	.size	internal_ping, .-internal_ping
	.section	.rodata.str1.8
	.align 8
.LC80:
	.string	"libpq is incorrectly linked to backend functions\n"
	.text
	.p2align 4,,15
	.type	connectDBStart, @function
connectDBStart:
.LFB933:
	.cfi_startproc
	xorl	%eax, %eax
	testq	%rdi, %rdi
	je	.L1020
	pushq	%rbx
	.cfi_def_cfa_offset 16
	.cfi_offset 3, -16
	cmpb	$0, 366(%rdi)
	movq	%rdi, %rbx
	jne	.L1023
.L1012:
	movl	$1, %esi
	movq	%rbx, %rdi
	call	pqDropConnection@PLT
	movl	$1, 336(%rbx)
	xorl	%eax, %eax
.L1009:
	popq	%rbx
	.cfi_remember_state
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L1023:
	.cfi_restore_state
	call	pg_link_canary_is_frontend@PLT
	testb	%al, %al
	leaq	928(%rbx), %rdi
	jne	.L1013
	leaq	.LC80(%rip), %rsi
	call	printfPQExpBuffer@PLT
	jmp	.L1012
	.p2align 4,,10
	.p2align 3
.L1020:
	.cfi_def_cfa_offset 8
	.cfi_restore 3
	rep ret
	.p2align 4,,10
	.p2align 3
.L1013:
	.cfi_def_cfa_offset 16
	.cfi_offset 3, -16
	movq	$0, 820(%rbx)
	movl	$0, 828(%rbx)
	movl	$0, 844(%rbx)
	call	resetPQExpBuffer@PLT
	movl	$-1, 396(%rbx)
	movb	$0, 720(%rbx)
	movq	%rbx, %rdi
	movb	$1, 721(%rbx)
	movl	$8, 336(%rbx)
	call	PQconnectPoll@PLT
	movl	%eax, %edx
	movl	$1, %eax
	cmpl	$2, %edx
	jne	.L1012
	jmp	.L1009
	.cfi_endproc
.LFE933:
	.size	connectDBStart, .-connectDBStart
	.section	.rodata.str1.8
	.align 8
.LC81:
	.string	"PGEventProc \"%s\" failed during PGEVT_CONNRESET event\n"
	.text
	.p2align 4,,15
	.globl	PQreset
	.type	PQreset, @function
PQreset:
.LFB945:
	.cfi_startproc
	pushq	%r13
	.cfi_def_cfa_offset 16
	.cfi_offset 13, -16
	pushq	%r12
	.cfi_def_cfa_offset 24
	.cfi_offset 12, -24
	pushq	%rbp
	.cfi_def_cfa_offset 32
	.cfi_offset 6, -32
	pushq	%rbx
	.cfi_def_cfa_offset 40
	.cfi_offset 3, -40
	subq	$24, %rsp
	.cfi_def_cfa_offset 64
	movq	%fs:40, %rax
	movq	%rax, 8(%rsp)
	xorl	%eax, %eax
	testq	%rdi, %rdi
	je	.L1024
	movq	%rdi, %rbx
	call	closePGconn
	movq	%rbx, %rdi
	call	connectDBStart
	testl	%eax, %eax
	jne	.L1041
.L1024:
	movq	8(%rsp), %rax
	xorq	%fs:40, %rax
	jne	.L1042
	addq	$24, %rsp
	.cfi_remember_state
	.cfi_def_cfa_offset 40
	popq	%rbx
	.cfi_def_cfa_offset 32
	popq	%rbp
	.cfi_def_cfa_offset 24
	popq	%r12
	.cfi_def_cfa_offset 16
	popq	%r13
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L1041:
	.cfi_restore_state
	movq	%rbx, %rdi
	call	connectDBComplete
	testl	%eax, %eax
	je	.L1024
	movl	328(%rbx), %eax
	testl	%eax, %eax
	jle	.L1024
	xorl	%r12d, %r12d
	movq	%rsp, %r13
	jmp	.L1028
	.p2align 4,,10
	.p2align 3
.L1026:
	addl	$1, %r12d
	cmpl	%r12d, 328(%rbx)
	jle	.L1024
.L1028:
	movslq	%r12d, %rax
	movq	%rbx, (%rsp)
	movq	%r13, %rsi
	leaq	(%rax,%rax,4), %rbp
	movq	320(%rbx), %rax
	movl	$1, %edi
	salq	$3, %rbp
	addq	%rbp, %rax
	movq	16(%rax), %rdx
	call	*(%rax)
	testl	%eax, %eax
	jne	.L1026
	movq	320(%rbx), %rax
	movl	$1, 336(%rbx)
	leaq	928(%rbx), %rdi
	leaq	.LC81(%rip), %rsi
	movq	8(%rax,%rbp), %rdx
	xorl	%eax, %eax
	call	printfPQExpBuffer@PLT
	jmp	.L1024
.L1042:
	call	__stack_chk_fail@PLT
	.cfi_endproc
.LFE945:
	.size	PQreset, .-PQreset
	.p2align 4,,15
	.globl	PQresetStart
	.type	PQresetStart, @function
PQresetStart:
.LFB946:
	.cfi_startproc
	testq	%rdi, %rdi
	je	.L1044
	pushq	%rbx
	.cfi_def_cfa_offset 16
	.cfi_offset 3, -16
	movq	%rdi, %rbx
	call	closePGconn
	movq	%rbx, %rdi
	popq	%rbx
	.cfi_restore 3
	.cfi_def_cfa_offset 8
	jmp	connectDBStart
	.p2align 4,,10
	.p2align 3
.L1044:
	xorl	%eax, %eax
	ret
	.cfi_endproc
.LFE946:
	.size	PQresetStart, .-PQresetStart
	.p2align 4,,15
	.globl	PQresetPoll
	.type	PQresetPoll, @function
PQresetPoll:
.LFB947:
	.cfi_startproc
	pushq	%r14
	.cfi_def_cfa_offset 16
	.cfi_offset 14, -16
	pushq	%r13
	.cfi_def_cfa_offset 24
	.cfi_offset 13, -24
	xorl	%r14d, %r14d
	pushq	%r12
	.cfi_def_cfa_offset 32
	.cfi_offset 12, -32
	pushq	%rbp
	.cfi_def_cfa_offset 40
	.cfi_offset 6, -40
	pushq	%rbx
	.cfi_def_cfa_offset 48
	.cfi_offset 3, -48
	subq	$16, %rsp
	.cfi_def_cfa_offset 64
	movq	%fs:40, %rax
	movq	%rax, 8(%rsp)
	xorl	%eax, %eax
	testq	%rdi, %rdi
	je	.L1048
	movq	%rdi, %rbx
	call	PQconnectPoll@PLT
	cmpl	$3, %eax
	movl	%eax, %r14d
	je	.L1060
.L1048:
	movq	8(%rsp), %rcx
	xorq	%fs:40, %rcx
	movl	%r14d, %eax
	jne	.L1061
	addq	$16, %rsp
	.cfi_remember_state
	.cfi_def_cfa_offset 48
	popq	%rbx
	.cfi_def_cfa_offset 40
	popq	%rbp
	.cfi_def_cfa_offset 32
	popq	%r12
	.cfi_def_cfa_offset 24
	popq	%r13
	.cfi_def_cfa_offset 16
	popq	%r14
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L1060:
	.cfi_restore_state
	movl	328(%rbx), %eax
	xorl	%r12d, %r12d
	movq	%rsp, %r13
	testl	%eax, %eax
	jg	.L1050
	jmp	.L1048
	.p2align 4,,10
	.p2align 3
.L1051:
	addl	$1, %r12d
	cmpl	%r12d, 328(%rbx)
	jle	.L1048
.L1050:
	movslq	%r12d, %rax
	movq	%rbx, (%rsp)
	movq	%r13, %rsi
	leaq	(%rax,%rax,4), %rbp
	movq	320(%rbx), %rax
	movl	$1, %edi
	salq	$3, %rbp
	addq	%rbp, %rax
	movq	16(%rax), %rdx
	call	*(%rax)
	testl	%eax, %eax
	jne	.L1051
	movq	320(%rbx), %rax
	movl	$1, 336(%rbx)
	leaq	928(%rbx), %rdi
	leaq	.LC81(%rip), %rsi
	xorl	%r14d, %r14d
	movq	8(%rax,%rbp), %rdx
	xorl	%eax, %eax
	call	printfPQExpBuffer@PLT
	jmp	.L1048
.L1061:
	call	__stack_chk_fail@PLT
	.cfi_endproc
.LFE947:
	.size	PQresetPoll, .-PQresetPoll
	.p2align 4,,15
	.globl	PQconninfo
	.type	PQconninfo, @function
PQconninfo:
.LFB972:
	.cfi_startproc
	pushq	%r13
	.cfi_def_cfa_offset 16
	.cfi_offset 13, -16
	pushq	%r12
	.cfi_def_cfa_offset 24
	.cfi_offset 12, -24
	pushq	%rbp
	.cfi_def_cfa_offset 32
	.cfi_offset 6, -32
	pushq	%rbx
	.cfi_def_cfa_offset 40
	.cfi_offset 3, -40
	subq	$40, %rsp
	.cfi_def_cfa_offset 80
	movq	%fs:40, %rax
	movq	%rax, 24(%rsp)
	xorl	%eax, %eax
	testq	%rdi, %rdi
	je	.L1065
	movq	%rsp, %rbp
	movq	%rdi, %r12
	movq	%rbp, %rdi
	call	initPQExpBuffer@PLT
	cmpq	$0, 16(%rsp)
	je	.L1065
	movq	%rbp, %rdi
	leaq	PQconninfoOptions(%rip), %rbx
	call	conninfo_init
	testq	%rax, %rax
	movq	%rax, %r13
	leaq	.LC2(%rip), %rsi
	jne	.L1066
	jmp	.L1069
	.p2align 4,,10
	.p2align 3
.L1068:
	movq	56(%rbx), %rax
	testq	%rax, %rax
	js	.L1067
	movq	(%r12,%rax), %rdx
	testq	%rdx, %rdx
	je	.L1067
	xorl	%r9d, %r9d
	movl	$1, %r8d
	movq	%rbp, %rcx
	movq	%r13, %rdi
	call	conninfo_storeval
.L1067:
	movq	64(%rbx), %rsi
.L1066:
	addq	$64, %rbx
	testq	%rsi, %rsi
	jne	.L1068
.L1069:
	movq	%rbp, %rdi
	call	termPQExpBuffer@PLT
.L1062:
	movq	24(%rsp), %rcx
	xorq	%fs:40, %rcx
	movq	%r13, %rax
	jne	.L1077
	addq	$40, %rsp
	.cfi_remember_state
	.cfi_def_cfa_offset 40
	popq	%rbx
	.cfi_def_cfa_offset 32
	popq	%rbp
	.cfi_def_cfa_offset 24
	popq	%r12
	.cfi_def_cfa_offset 16
	popq	%r13
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L1065:
	.cfi_restore_state
	xorl	%r13d, %r13d
	jmp	.L1062
.L1077:
	call	__stack_chk_fail@PLT
	.cfi_endproc
.LFE972:
	.size	PQconninfo, .-PQconninfo
	.p2align 4,,15
	.globl	PQconninfoFree
	.type	PQconninfoFree, @function
PQconninfoFree:
.LFB973:
	.cfi_startproc
	testq	%rdi, %rdi
	je	.L1078
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset 6, -16
	pushq	%rbx
	.cfi_def_cfa_offset 24
	.cfi_offset 3, -24
	movq	%rdi, %rbp
	subq	$8, %rsp
	.cfi_def_cfa_offset 32
	cmpq	$0, (%rdi)
	je	.L1080
	movq	%rdi, %rbx
	.p2align 4,,10
	.p2align 3
.L1082:
	movq	24(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L1081
	call	free@PLT
.L1081:
	addq	$56, %rbx
	cmpq	$0, (%rbx)
	jne	.L1082
.L1080:
	addq	$8, %rsp
	.cfi_def_cfa_offset 24
	movq	%rbp, %rdi
	popq	%rbx
	.cfi_restore 3
	.cfi_def_cfa_offset 16
	popq	%rbp
	.cfi_restore 6
	.cfi_def_cfa_offset 8
	jmp	free@PLT
	.p2align 4,,10
	.p2align 3
.L1078:
	rep ret
	.cfi_endproc
.LFE973:
	.size	PQconninfoFree, .-PQconninfoFree
	.p2align 4,,15
	.globl	PQdb
	.type	PQdb, @function
PQdb:
.LFB974:
	.cfi_startproc
	testq	%rdi, %rdi
	je	.L1092
	movq	88(%rdi), %rax
	ret
	.p2align 4,,10
	.p2align 3
.L1092:
	xorl	%eax, %eax
	ret
	.cfi_endproc
.LFE974:
	.size	PQdb, .-PQdb
	.p2align 4,,15
	.globl	PQuser
	.type	PQuser, @function
PQuser:
.LFB975:
	.cfi_startproc
	testq	%rdi, %rdi
	je	.L1095
	movq	104(%rdi), %rax
	ret
	.p2align 4,,10
	.p2align 3
.L1095:
	xorl	%eax, %eax
	ret
	.cfi_endproc
.LFE975:
	.size	PQuser, .-PQuser
	.section	.rodata.str1.1
.LC82:
	.string	""
	.text
	.p2align 4,,15
	.globl	PQpass
	.type	PQpass, @function
PQpass:
.LFB976:
	.cfi_startproc
	testq	%rdi, %rdi
	je	.L1099
	movq	400(%rdi), %rax
	testq	%rax, %rax
	je	.L1098
	movslq	396(%rdi), %rdx
	leaq	(%rdx,%rdx,4), %rdx
	movq	32(%rax,%rdx,8), %rax
	testq	%rax, %rax
	je	.L1098
	rep ret
	.p2align 4,,10
	.p2align 3
.L1098:
	movq	112(%rdi), %rax
	leaq	.LC82(%rip), %rdx
	testq	%rax, %rax
	cmove	%rdx, %rax
	ret
	.p2align 4,,10
	.p2align 3
.L1099:
	xorl	%eax, %eax
	ret
	.cfi_endproc
.LFE976:
	.size	PQpass, .-PQpass
	.p2align 4,,15
	.globl	PQhost
	.type	PQhost, @function
PQhost:
.LFB977:
	.cfi_startproc
	testq	%rdi, %rdi
	je	.L1106
	movq	400(%rdi), %rdx
	leaq	.LC82(%rip), %rax
	testq	%rdx, %rdx
	je	.L1103
	movslq	396(%rdi), %rax
	leaq	(%rax,%rax,4), %rax
	leaq	(%rdx,%rax,8), %rdx
	movq	8(%rdx), %rax
	testq	%rax, %rax
	je	.L1105
	cmpb	$0, (%rax)
	jne	.L1103
.L1105:
	movq	16(%rdx), %rax
	testq	%rax, %rax
	je	.L1108
	cmpb	$0, (%rax)
	leaq	.LC82(%rip), %rdx
	cmove	%rdx, %rax
	ret
	.p2align 4,,10
	.p2align 3
.L1108:
	leaq	.LC82(%rip), %rax
.L1103:
	rep ret
	.p2align 4,,10
	.p2align 3
.L1106:
	xorl	%eax, %eax
	ret
	.cfi_endproc
.LFE977:
	.size	PQhost, .-PQhost
	.p2align 4,,15
	.globl	PQhostaddr
	.type	PQhostaddr, @function
PQhostaddr:
.LFB978:
	.cfi_startproc
	testq	%rdi, %rdi
	je	.L1116
	cmpq	$0, 400(%rdi)
	leaq	.LC82(%rip), %rax
	je	.L1114
	movq	408(%rdi), %rax
	leaq	.LC82(%rip), %rdx
	testq	%rax, %rax
	cmove	%rdx, %rax
	ret
	.p2align 4,,10
	.p2align 3
.L1116:
	xorl	%eax, %eax
.L1114:
	rep ret
	.cfi_endproc
.LFE978:
	.size	PQhostaddr, .-PQhostaddr
	.p2align 4,,15
	.globl	PQport
	.type	PQport, @function
PQport:
.LFB979:
	.cfi_startproc
	testq	%rdi, %rdi
	je	.L1120
	movq	400(%rdi), %rdx
	leaq	.LC82(%rip), %rax
	testq	%rdx, %rdx
	je	.L1118
	movslq	396(%rdi), %rax
	leaq	(%rax,%rax,4), %rax
	movq	24(%rdx,%rax,8), %rax
	ret
	.p2align 4,,10
	.p2align 3
.L1120:
	xorl	%eax, %eax
.L1118:
	rep ret
	.cfi_endproc
.LFE979:
	.size	PQport, .-PQport
	.p2align 4,,15
	.globl	PQtty
	.type	PQtty, @function
PQtty:
.LFB980:
	.cfi_startproc
	testq	%rdi, %rdi
	je	.L1125
	movq	32(%rdi), %rax
	ret
	.p2align 4,,10
	.p2align 3
.L1125:
	xorl	%eax, %eax
	ret
	.cfi_endproc
.LFE980:
	.size	PQtty, .-PQtty
	.p2align 4,,15
	.globl	PQoptions
	.type	PQoptions, @function
PQoptions:
.LFB981:
	.cfi_startproc
	testq	%rdi, %rdi
	je	.L1128
	movq	64(%rdi), %rax
	ret
	.p2align 4,,10
	.p2align 3
.L1128:
	xorl	%eax, %eax
	ret
	.cfi_endproc
.LFE981:
	.size	PQoptions, .-PQoptions
	.p2align 4,,15
	.globl	PQstatus
	.type	PQstatus, @function
PQstatus:
.LFB982:
	.cfi_startproc
	testq	%rdi, %rdi
	movl	$1, %eax
	je	.L1129
	movl	336(%rdi), %eax
.L1129:
	rep ret
	.cfi_endproc
.LFE982:
	.size	PQstatus, .-PQstatus
	.p2align 4,,15
	.globl	PQtransactionStatus
	.type	PQtransactionStatus, @function
PQtransactionStatus:
.LFB983:
	.cfi_startproc
	testq	%rdi, %rdi
	movl	$4, %eax
	je	.L1133
	movl	336(%rdi), %ecx
	testl	%ecx, %ecx
	jne	.L1133
	movl	340(%rdi), %edx
	movl	$1, %eax
	testl	%edx, %edx
	jne	.L1133
	movl	344(%rdi), %eax
	ret
	.p2align 4,,10
	.p2align 3
.L1133:
	rep ret
	.cfi_endproc
.LFE983:
	.size	PQtransactionStatus, .-PQtransactionStatus
	.p2align 4,,15
	.globl	PQparameterStatus
	.type	PQparameterStatus, @function
PQparameterStatus:
.LFB984:
	.cfi_startproc
	testq	%rdi, %rdi
	je	.L1142
	testq	%rsi, %rsi
	je	.L1142
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset 6, -16
	pushq	%rbx
	.cfi_def_cfa_offset 24
	.cfi_offset 3, -24
	movq	%rsi, %rbp
	subq	$8, %rsp
	.cfi_def_cfa_offset 32
	movq	776(%rdi), %rbx
	testq	%rbx, %rbx
	jne	.L1141
	jmp	.L1143
	.p2align 4,,10
	.p2align 3
.L1140:
	movq	(%rbx), %rbx
	testq	%rbx, %rbx
	je	.L1143
.L1141:
	movq	8(%rbx), %rdi
	movq	%rbp, %rsi
	call	strcmp@PLT
	testl	%eax, %eax
	jne	.L1140
	movq	16(%rbx), %rax
	addq	$8, %rsp
	.cfi_remember_state
	.cfi_def_cfa_offset 24
	popq	%rbx
	.cfi_def_cfa_offset 16
	popq	%rbp
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L1143:
	.cfi_restore_state
	addq	$8, %rsp
	.cfi_def_cfa_offset 24
	xorl	%eax, %eax
	popq	%rbx
	.cfi_def_cfa_offset 16
	popq	%rbp
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L1142:
	.cfi_restore 3
	.cfi_restore 6
	xorl	%eax, %eax
	ret
	.cfi_endproc
.LFE984:
	.size	PQparameterStatus, .-PQparameterStatus
	.p2align 4,,15
	.globl	PQprotocolVersion
	.type	PQprotocolVersion, @function
PQprotocolVersion:
.LFB985:
	.cfi_startproc
	xorl	%eax, %eax
	testq	%rdi, %rdi
	je	.L1149
	cmpl	$1, 336(%rdi)
	je	.L1149
	movzwl	698(%rdi), %eax
.L1149:
	rep ret
	.cfi_endproc
.LFE985:
	.size	PQprotocolVersion, .-PQprotocolVersion
	.p2align 4,,15
	.globl	PQserverVersion
	.type	PQserverVersion, @function
PQserverVersion:
.LFB986:
	.cfi_startproc
	xorl	%eax, %eax
	testq	%rdi, %rdi
	je	.L1154
	cmpl	$1, 336(%rdi)
	je	.L1154
	movl	700(%rdi), %eax
.L1154:
	rep ret
	.cfi_endproc
.LFE986:
	.size	PQserverVersion, .-PQserverVersion
	.section	.rodata.str1.1
.LC83:
	.string	"connection pointer is NULL\n"
	.text
	.p2align 4,,15
	.globl	PQerrorMessage
	.type	PQerrorMessage, @function
PQerrorMessage:
.LFB987:
	.cfi_startproc
	testq	%rdi, %rdi
	leaq	.LC83(%rip), %rax
	je	.L1159
	movq	928(%rdi), %rax
.L1159:
	rep ret
	.cfi_endproc
.LFE987:
	.size	PQerrorMessage, .-PQerrorMessage
	.p2align 4,,15
	.globl	PQsocket
	.type	PQsocket, @function
PQsocket:
.LFB988:
	.cfi_startproc
	testq	%rdi, %rdi
	je	.L1165
	movl	416(%rdi), %eax
	ret
.L1165:
	movl	$-1, %eax
	ret
	.cfi_endproc
.LFE988:
	.size	PQsocket, .-PQsocket
	.p2align 4,,15
	.globl	PQbackendPID
	.type	PQbackendPID, @function
PQbackendPID:
.LFB989:
	.cfi_startproc
	xorl	%eax, %eax
	testq	%rdi, %rdi
	je	.L1166
	movl	336(%rdi), %edx
	testl	%edx, %edx
	jne	.L1166
	movl	764(%rdi), %eax
.L1166:
	rep ret
	.cfi_endproc
.LFE989:
	.size	PQbackendPID, .-PQbackendPID
	.p2align 4,,15
	.globl	PQconnectionNeedsPassword
	.type	PQconnectionNeedsPassword, @function
PQconnectionNeedsPassword:
.LFB990:
	.cfi_startproc
	testq	%rdi, %rdi
	je	.L1181
	pushq	%rbx
	.cfi_def_cfa_offset 16
	.cfi_offset 3, -16
	movq	%rdi, %rbx
	call	PQpass@PLT
	cmpb	$0, 705(%rbx)
	movq	%rax, %rdx
	je	.L1172
	testq	%rdx, %rdx
	movl	$1, %eax
	je	.L1171
	cmpb	$0, (%rdx)
	je	.L1171
.L1172:
	xorl	%eax, %eax
.L1171:
	popq	%rbx
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L1181:
	.cfi_restore 3
	xorl	%eax, %eax
	ret
	.cfi_endproc
.LFE990:
	.size	PQconnectionNeedsPassword, .-PQconnectionNeedsPassword
	.p2align 4,,15
	.globl	PQconnectionUsedPassword
	.type	PQconnectionUsedPassword, @function
PQconnectionUsedPassword:
.LFB991:
	.cfi_startproc
	xorl	%eax, %eax
	testq	%rdi, %rdi
	je	.L1185
	movzbl	705(%rdi), %eax
.L1185:
	rep ret
	.cfi_endproc
.LFE991:
	.size	PQconnectionUsedPassword, .-PQconnectionUsedPassword
	.p2align 4,,15
	.globl	PQclientEncoding
	.type	PQclientEncoding, @function
PQclientEncoding:
.LFB992:
	.cfi_startproc
	testq	%rdi, %rdi
	je	.L1192
	movl	336(%rdi), %eax
	testl	%eax, %eax
	jne	.L1192
	movl	784(%rdi), %eax
	ret
	.p2align 4,,10
	.p2align 3
.L1192:
	movl	$-1, %eax
	ret
	.cfi_endproc
.LFE992:
	.size	PQclientEncoding, .-PQclientEncoding
	.section	.rodata.str1.1
.LC84:
	.string	"auto"
.LC85:
	.string	"client_encoding"
	.text
	.p2align 4,,15
	.globl	PQsetClientEncoding
	.type	PQsetClientEncoding, @function
PQsetClientEncoding:
.LFB993:
	.cfi_startproc
	pushq	%r13
	.cfi_def_cfa_offset 16
	.cfi_offset 13, -16
	pushq	%r12
	.cfi_def_cfa_offset 24
	.cfi_offset 12, -24
	pushq	%rbp
	.cfi_def_cfa_offset 32
	.cfi_offset 6, -32
	pushq	%rbx
	.cfi_def_cfa_offset 40
	.cfi_offset 3, -40
	subq	$168, %rsp
	.cfi_def_cfa_offset 208
	movq	%fs:40, %rax
	movq	%rax, 152(%rsp)
	xorl	%eax, %eax
	testq	%rdi, %rdi
	je	.L1196
	movl	336(%rdi), %eax
	movq	%rdi, %rbp
	testl	%eax, %eax
	jne	.L1196
	testq	%rsi, %rsi
	movq	%rsi, %rbx
	je	.L1196
	leaq	.LC84(%rip), %rdi
	movl	$5, %ecx
	repz cmpsb
	seta	%al
	sbbb	$0, %al
	testb	%al, %al
	je	.L1206
.L1197:
	movq	%rbx, %rdi
	call	strlen@PLT
	addq	$28, %rax
	cmpq	$128, %rax
	ja	.L1196
	leaq	16(%rsp), %r12
	leaq	query.25836(%rip), %rsi
	movq	%rbx, %rdx
	xorl	%eax, %eax
	movq	%r12, %rdi
	call	pg_sprintf@PLT
	movq	%r12, %rsi
	movq	%rbp, %rdi
	call	PQexec@PLT
	testq	%rax, %rax
	je	.L1196
	cmpl	$1, 40(%rax)
	movl	$-1, %r13d
	je	.L1207
.L1198:
	movq	%rax, %rdi
	call	PQclear@PLT
.L1193:
	movq	152(%rsp), %rdx
	xorq	%fs:40, %rdx
	movl	%r13d, %eax
	jne	.L1208
	addq	$168, %rsp
	.cfi_remember_state
	.cfi_def_cfa_offset 40
	popq	%rbx
	.cfi_def_cfa_offset 32
	popq	%rbp
	.cfi_def_cfa_offset 24
	popq	%r12
	.cfi_def_cfa_offset 16
	popq	%r13
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L1206:
	.cfi_restore_state
	movl	$1, %esi
	xorl	%edi, %edi
	call	pg_get_encoding_from_locale@PLT
	movl	%eax, %edi
	call	pg_encoding_to_char@PLT
	movq	%rax, %rbx
	jmp	.L1197
	.p2align 4,,10
	.p2align 3
.L1207:
	xorl	%r13d, %r13d
	cmpw	$2, 698(%rbp)
	ja	.L1198
	leaq	.LC85(%rip), %rsi
	movq	%rbx, %rdx
	movq	%rbp, %rdi
	movq	%rax, 8(%rsp)
	call	pqSaveParameterStatus@PLT
	movq	8(%rsp), %rax
	jmp	.L1198
	.p2align 4,,10
	.p2align 3
.L1196:
	movl	$-1, %r13d
	jmp	.L1193
.L1208:
	call	__stack_chk_fail@PLT
	.cfi_endproc
.LFE993:
	.size	PQsetClientEncoding, .-PQsetClientEncoding
	.p2align 4,,15
	.globl	PQsetErrorVerbosity
	.type	PQsetErrorVerbosity, @function
PQsetErrorVerbosity:
.LFB994:
	.cfi_startproc
	testq	%rdi, %rdi
	movl	$1, %eax
	je	.L1209
	movl	792(%rdi), %eax
	movl	%esi, 792(%rdi)
.L1209:
	rep ret
	.cfi_endproc
.LFE994:
	.size	PQsetErrorVerbosity, .-PQsetErrorVerbosity
	.p2align 4,,15
	.globl	PQsetErrorContextVisibility
	.type	PQsetErrorContextVisibility, @function
PQsetErrorContextVisibility:
.LFB995:
	.cfi_startproc
	testq	%rdi, %rdi
	movl	$1, %eax
	je	.L1213
	movl	796(%rdi), %eax
	movl	%esi, 796(%rdi)
.L1213:
	rep ret
	.cfi_endproc
.LFE995:
	.size	PQsetErrorContextVisibility, .-PQsetErrorContextVisibility
	.p2align 4,,15
	.globl	PQuntrace
	.type	PQuntrace, @function
PQuntrace:
.LFB997:
	.cfi_startproc
	testq	%rdi, %rdi
	je	.L1226
	movq	280(%rdi), %rax
	testq	%rax, %rax
	je	.L1226
	pushq	%rbx
	.cfi_def_cfa_offset 16
	.cfi_offset 3, -16
	movq	%rdi, %rbx
	movq	%rax, %rdi
	call	fflush@PLT
	movq	$0, 280(%rbx)
	popq	%rbx
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L1226:
	.cfi_restore 3
	rep ret
	.cfi_endproc
.LFE997:
	.size	PQuntrace, .-PQuntrace
	.p2align 4,,15
	.globl	PQtrace
	.type	PQtrace, @function
PQtrace:
.LFB996:
	.cfi_startproc
	testq	%rdi, %rdi
	je	.L1235
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset 6, -16
	pushq	%rbx
	.cfi_def_cfa_offset 24
	.cfi_offset 3, -24
	movq	%rsi, %rbp
	movq	%rdi, %rbx
	subq	$8, %rsp
	.cfi_def_cfa_offset 32
	call	PQuntrace@PLT
	movq	%rbp, 280(%rbx)
	addq	$8, %rsp
	.cfi_def_cfa_offset 24
	popq	%rbx
	.cfi_def_cfa_offset 16
	popq	%rbp
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L1235:
	.cfi_restore 3
	.cfi_restore 6
	rep ret
	.cfi_endproc
.LFE996:
	.size	PQtrace, .-PQtrace
	.p2align 4,,15
	.globl	PQsetNoticeReceiver
	.type	PQsetNoticeReceiver, @function
PQsetNoticeReceiver:
.LFB998:
	.cfi_startproc
	testq	%rdi, %rdi
	je	.L1240
	testq	%rsi, %rsi
	movq	288(%rdi), %rax
	je	.L1238
	movq	%rsi, 288(%rdi)
	movq	%rdx, 296(%rdi)
	ret
	.p2align 4,,10
	.p2align 3
.L1240:
	xorl	%eax, %eax
.L1238:
	rep ret
	.cfi_endproc
.LFE998:
	.size	PQsetNoticeReceiver, .-PQsetNoticeReceiver
	.p2align 4,,15
	.globl	PQsetNoticeProcessor
	.type	PQsetNoticeProcessor, @function
PQsetNoticeProcessor:
.LFB999:
	.cfi_startproc
	testq	%rdi, %rdi
	je	.L1246
	testq	%rsi, %rsi
	movq	304(%rdi), %rax
	je	.L1244
	movq	%rsi, 304(%rdi)
	movq	%rdx, 312(%rdi)
	ret
	.p2align 4,,10
	.p2align 3
.L1246:
	xorl	%eax, %eax
.L1244:
	rep ret
	.cfi_endproc
.LFE999:
	.size	PQsetNoticeProcessor, .-PQsetNoticeProcessor
	.p2align 4,,15
	.globl	pqGetHomeDirectory
	.type	pqGetHomeDirectory, @function
pqGetHomeDirectory:
.LFB1007:
	.cfi_startproc
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset 6, -16
	pushq	%rbx
	.cfi_def_cfa_offset 24
	.cfi_offset 3, -24
	movq	%rdi, %rbp
	movl	%esi, %ebx
	subq	$8280, %rsp
	.cfi_def_cfa_offset 8304
	movq	%fs:40, %rax
	movq	%rax, 8264(%rsp)
	xorl	%eax, %eax
	movq	$0, 8(%rsp)
	call	geteuid@PLT
	leaq	64(%rsp), %rdx
	leaq	16(%rsp), %rsi
	leaq	8(%rsp), %r8
	movl	$8192, %ecx
	movl	%eax, %edi
	call	pqGetpwuid@PLT
	movq	8(%rsp), %rcx
	xorl	%eax, %eax
	testq	%rcx, %rcx
	je	.L1250
	movq	32(%rcx), %rsi
	movslq	%ebx, %rdx
	movq	%rbp, %rdi
	call	strlcpy@PLT
	movl	$1, %eax
.L1250:
	movq	8264(%rsp), %rbx
	xorq	%fs:40, %rbx
	jne	.L1257
	addq	$8280, %rsp
	.cfi_remember_state
	.cfi_def_cfa_offset 24
	popq	%rbx
	.cfi_def_cfa_offset 16
	popq	%rbp
	.cfi_def_cfa_offset 8
	ret
.L1257:
	.cfi_restore_state
	call	__stack_chk_fail@PLT
	.cfi_endproc
.LFE1007:
	.size	pqGetHomeDirectory, .-pqGetHomeDirectory
	.section	.rodata.str1.1
.LC86:
	.string	"/usr/local/pgsql/etc"
.LC87:
	.string	"PGSERVICE"
.LC88:
	.string	"PGSERVICEFILE"
.LC89:
	.string	".pg_service.conf"
.LC90:
	.string	"%s/%s"
.LC91:
	.string	"PGSYSCONFDIR"
.LC92:
	.string	"%s/pg_service.conf"
	.section	.rodata.str1.8
	.align 8
.LC93:
	.string	"definition of service \"%s\" not found\n"
	.text
	.p2align 4,,15
	.type	parseServiceInfo, @function
parseServiceInfo:
.LFB954:
	.cfi_startproc
	pushq	%r14
	.cfi_def_cfa_offset 16
	.cfi_offset 14, -16
	pushq	%r13
	.cfi_def_cfa_offset 24
	.cfi_offset 13, -24
	movq	%rsi, %r13
	pushq	%r12
	.cfi_def_cfa_offset 32
	.cfi_offset 12, -32
	pushq	%rbp
	.cfi_def_cfa_offset 40
	.cfi_offset 6, -40
	leaq	.LC2(%rip), %rsi
	pushq	%rbx
	.cfi_def_cfa_offset 48
	.cfi_offset 3, -48
	movq	%rdi, %rbp
	subq	$2224, %rsp
	.cfi_def_cfa_offset 2272
	movq	%fs:40, %rax
	movq	%rax, 2216(%rsp)
	xorl	%eax, %eax
	call	conninfo_getval
	testq	%rax, %rax
	movb	$0, 15(%rsp)
	movq	%rax, %rbx
	je	.L1280
.L1259:
	leaq	.LC88(%rip), %rdi
	call	getenv@PLT
	testq	%rax, %rax
	je	.L1261
	leaq	160(%rsp), %r12
	movl	$1024, %edx
	movq	%rax, %rsi
	movq	%r12, %rdi
	call	strlcpy@PLT
.L1262:
	leaq	15(%rsp), %r8
	movq	%r13, %rcx
	movq	%rbp, %rdx
	movq	%rbx, %rsi
	movq	%r12, %rdi
	call	parseServiceFile
	cmpb	$0, 15(%rsp)
	jne	.L1258
	testl	%eax, %eax
	je	.L1279
.L1258:
	movq	2216(%rsp), %rcx
	xorq	%fs:40, %rcx
	jne	.L1281
	addq	$2224, %rsp
	.cfi_remember_state
	.cfi_def_cfa_offset 48
	popq	%rbx
	.cfi_def_cfa_offset 40
	popq	%rbp
	.cfi_def_cfa_offset 32
	popq	%r12
	.cfi_def_cfa_offset 24
	popq	%r13
	.cfi_def_cfa_offset 16
	popq	%r14
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L1279:
	.cfi_restore_state
	leaq	16(%rsp), %r14
.L1264:
	leaq	.LC91(%rip), %rdi
	call	getenv@PLT
	testq	%rax, %rax
	leaq	.LC86(%rip), %rcx
	je	.L1265
	leaq	.LC91(%rip), %rdi
	call	getenv@PLT
	movq	%rax, %rcx
.L1265:
	leaq	.LC92(%rip), %rdx
	movl	$1024, %esi
	movq	%r12, %rdi
	xorl	%eax, %eax
	call	pg_snprintf@PLT
	movq	%r14, %rdx
	movq	%r12, %rsi
	movl	$1, %edi
	call	__xstat@PLT
	testl	%eax, %eax
	je	.L1266
.L1269:
	cmpb	$0, 15(%rsp)
	je	.L1282
	xorl	%eax, %eax
	jmp	.L1258
	.p2align 4,,10
	.p2align 3
.L1261:
	leaq	1184(%rsp), %r14
	movl	$1024, %esi
	leaq	160(%rsp), %r12
	movq	%r14, %rdi
	call	pqGetHomeDirectory@PLT
	testb	%al, %al
	je	.L1279
	movq	%r14, %rcx
	leaq	.LC89(%rip), %r8
	leaq	.LC90(%rip), %rdx
	leaq	16(%rsp), %r14
	movl	$1024, %esi
	movq	%r12, %rdi
	xorl	%eax, %eax
	call	pg_snprintf@PLT
	movq	%r14, %rdx
	movq	%r12, %rsi
	movl	$1, %edi
	call	__xstat@PLT
	testl	%eax, %eax
	jne	.L1264
	jmp	.L1262
	.p2align 4,,10
	.p2align 3
.L1266:
	leaq	15(%rsp), %r8
	movq	%r13, %rcx
	movq	%rbp, %rdx
	movq	%rbx, %rsi
	movq	%r12, %rdi
	call	parseServiceFile
	testl	%eax, %eax
	je	.L1269
	jmp	.L1258
	.p2align 4,,10
	.p2align 3
.L1280:
	leaq	.LC87(%rip), %rdi
	call	getenv@PLT
	testq	%rax, %rax
	movq	%rax, %rbx
	jne	.L1259
	xorl	%eax, %eax
	jmp	.L1258
	.p2align 4,,10
	.p2align 3
.L1282:
	leaq	.LC93(%rip), %rsi
	xorl	%eax, %eax
	movq	%rbx, %rdx
	movq	%r13, %rdi
	call	printfPQExpBuffer@PLT
	movl	$3, %eax
	jmp	.L1258
.L1281:
	call	__stack_chk_fail@PLT
	.cfi_endproc
.LFE954:
	.size	parseServiceInfo, .-parseServiceInfo
	.section	.rodata.str1.1
.LC94:
	.string	"PGREQUIRESSL"
.LC95:
	.string	"user"
	.text
	.p2align 4,,15
	.type	conninfo_add_defaults, @function
conninfo_add_defaults:
.LFB963:
	.cfi_startproc
	pushq	%r15
	.cfi_def_cfa_offset 16
	.cfi_offset 15, -16
	pushq	%r14
	.cfi_def_cfa_offset 24
	.cfi_offset 14, -24
	movq	%rdi, %r15
	pushq	%r13
	.cfi_def_cfa_offset 32
	.cfi_offset 13, -32
	pushq	%r12
	.cfi_def_cfa_offset 40
	.cfi_offset 12, -40
	pushq	%rbp
	.cfi_def_cfa_offset 48
	.cfi_offset 6, -48
	pushq	%rbx
	.cfi_def_cfa_offset 56
	.cfi_offset 3, -56
	movq	%rsi, %rbp
	subq	$8, %rsp
	.cfi_def_cfa_offset 64
	call	parseServiceInfo
	testl	%eax, %eax
	setne	%bl
	testq	%rbp, %rbp
	setne	%al
	andb	%al, %bl
	jne	.L1284
	cmpq	$0, (%r15)
	leaq	.LC15(%rip), %r12
	leaq	.LC95(%rip), %r13
	leaq	.LC94(%rip), %r14
	jne	.L1285
	jmp	.L1294
	.p2align 4,,10
	.p2align 3
.L1288:
	addq	$56, %r15
	cmpq	$0, (%r15)
	je	.L1294
.L1285:
	cmpq	$0, 24(%r15)
	jne	.L1288
	movq	8(%r15), %rdi
	testq	%rdi, %rdi
	je	.L1289
	call	getenv@PLT
	testq	%rax, %rax
	je	.L1289
	movq	%rax, %rdi
	jmp	.L1317
	.p2align 4,,10
	.p2align 3
.L1289:
	movq	(%r15), %rsi
	movl	$8, %ecx
	movq	%r12, %rdi
	repz cmpsb
	seta	%al
	sbbb	$0, %al
	testb	%al, %al
	je	.L1318
.L1291:
	movq	16(%r15), %rdi
	testq	%rdi, %rdi
	je	.L1293
.L1317:
	call	strdup@PLT
	testq	%rax, %rax
	movq	%rax, 24(%r15)
	jne	.L1288
.L1292:
	testq	%rbp, %rbp
	jne	.L1319
.L1284:
	xorl	%ebx, %ebx
.L1283:
	addq	$8, %rsp
	.cfi_remember_state
	.cfi_def_cfa_offset 56
	movl	%ebx, %eax
	popq	%rbx
	.cfi_def_cfa_offset 48
	popq	%rbp
	.cfi_def_cfa_offset 40
	popq	%r12
	.cfi_def_cfa_offset 32
	popq	%r13
	.cfi_def_cfa_offset 24
	popq	%r14
	.cfi_def_cfa_offset 16
	popq	%r15
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L1293:
	.cfi_restore_state
	movq	(%r15), %rsi
	movl	$5, %ecx
	movq	%r13, %rdi
	repz cmpsb
	seta	%al
	sbbb	$0, %al
	testb	%al, %al
	jne	.L1288
	xorl	%edi, %edi
	addq	$56, %r15
	call	pg_fe_getauthname@PLT
	movq	%rax, -32(%r15)
	cmpq	$0, (%r15)
	jne	.L1285
	.p2align 4,,10
	.p2align 3
.L1294:
	movl	$1, %ebx
	jmp	.L1283
	.p2align 4,,10
	.p2align 3
.L1318:
	movq	%r14, %rdi
	call	getenv@PLT
	testq	%rax, %rax
	je	.L1291
	cmpb	$49, (%rax)
	jne	.L1291
	leaq	.LC14(%rip), %rdi
	call	strdup@PLT
	testq	%rax, %rax
	movq	%rax, 24(%r15)
	jne	.L1288
	jmp	.L1292
	.p2align 4,,10
	.p2align 3
.L1319:
	leaq	.LC3(%rip), %rsi
	movq	%rbp, %rdi
	xorl	%eax, %eax
	call	printfPQExpBuffer@PLT
	jmp	.L1283
	.cfi_endproc
.LFE963:
	.size	conninfo_add_defaults, .-conninfo_add_defaults
	.section	.rodata.str1.8
	.align 8
.LC96:
	.string	"invalid URI propagated to internal parser routine: \"%s\"\n"
	.section	.rodata.str1.1
.LC97:
	.string	"password"
	.section	.rodata.str1.8
	.align 8
.LC98:
	.string	"end of string reached when looking for matching \"]\" in IPv6 host address in URI: \"%s\"\n"
	.align 8
.LC99:
	.string	"IPv6 host address may not be empty in URI: \"%s\"\n"
	.align 8
.LC100:
	.string	"unexpected character \"%c\" at position %d in URI (expected \":\" or \"/\"): \"%s\"\n"
	.section	.rodata.str1.1
.LC101:
	.string	"host"
.LC102:
	.string	"dbname"
	.section	.rodata.str1.8
	.align 8
.LC103:
	.string	"extra key/value separator \"=\" in URI query parameter: \"%s\"\n"
	.align 8
.LC104:
	.string	"missing key/value separator \"=\" in URI query parameter: \"%s\"\n"
	.section	.rodata.str1.1
.LC105:
	.string	"ssl"
.LC106:
	.string	"true"
	.section	.rodata.str1.8
	.align 8
.LC107:
	.string	"missing \"=\" after \"%s\" in connection info string\n"
	.align 8
.LC108:
	.string	"unterminated quoted string in connection info string\n"
	.align 8
.LC109:
	.string	"invalid URI query parameter: \"%s\"\n"
	.text
	.p2align 4,,15
	.type	parse_connection_string, @function
parse_connection_string:
.LFB958:
	.cfi_startproc
	pushq	%r15
	.cfi_def_cfa_offset 16
	.cfi_offset 15, -16
	pushq	%r14
	.cfi_def_cfa_offset 24
	.cfi_offset 14, -24
	pushq	%r13
	.cfi_def_cfa_offset 32
	.cfi_offset 13, -32
	pushq	%r12
	.cfi_def_cfa_offset 40
	.cfi_offset 12, -40
	movq	%rdi, %r13
	pushq	%rbp
	.cfi_def_cfa_offset 48
	.cfi_offset 6, -48
	pushq	%rbx
	.cfi_def_cfa_offset 56
	.cfi_offset 3, -56
	movq	%rsi, %rbp
	subq	$136, %rsp
	.cfi_def_cfa_offset 192
	movq	%fs:40, %rax
	movq	%rax, 120(%rsp)
	xorl	%eax, %eax
	movl	%edx, 8(%rsp)
	call	uri_prefix_length
	testl	%eax, %eax
	movq	%rbp, %rdi
	jne	.L1558
	call	conninfo_init
	testq	%rax, %rax
	movq	%rax, %rbx
	je	.L1549
	movq	%r13, %rdi
	call	strdup@PLT
	testq	%rax, %rax
	movq	%rax, %r12
	movq	%rax, %r13
	je	.L1559
.L1391:
	movzbl	0(%r13), %r14d
	testb	%r14b, %r14b
	je	.L1560
.L1421:
	call	__ctype_b_loc@PLT
	movq	(%rax), %rcx
	movzbl	%r14b, %eax
	testb	$32, 1(%rcx,%rax,2)
	jne	.L1392
	cmpb	$61, %r14b
	je	.L1443
	movq	%r13, %rdx
	jmp	.L1394
	.p2align 4,,10
	.p2align 3
.L1562:
	testb	%al, %al
	je	.L1397
	testb	$32, 1(%rcx,%rax,2)
	jne	.L1561
	movq	%rsi, %rdx
.L1394:
	movzbl	1(%rdx), %eax
	leaq	1(%rdx), %rsi
	cmpb	$61, %al
	jne	.L1562
.L1397:
	cmpb	$61, %al
	jne	.L1400
.L1393:
	movzbl	1(%rsi), %eax
	leaq	1(%rsi), %rdx
	movb	$0, (%rsi)
	testb	%al, %al
	jne	.L1406
	jmp	.L1403
	.p2align 4,,10
	.p2align 3
.L1563:
	addq	$1, %rdx
	movzbl	(%rdx), %eax
	testb	%al, %al
	je	.L1403
.L1406:
	movzbl	%al, %esi
	testb	$32, 1(%rcx,%rsi,2)
	jne	.L1563
	cmpb	$39, %al
	jne	.L1564
	movzbl	1(%rdx), %eax
	leaq	1(%rdx), %rcx
	movq	%rdx, %rsi
.L1415:
	testb	%al, %al
	je	.L1418
.L1416:
	cmpb	$92, %al
	leaq	1(%rcx), %r14
	je	.L1565
	cmpb	$39, %al
	je	.L1566
	movb	%al, (%rsi)
	movzbl	1(%rcx), %eax
	addq	$1, %rsi
	movq	%r14, %rcx
	testb	%al, %al
	jne	.L1416
.L1418:
	leaq	.LC108(%rip), %rsi
	movq	%rbp, %rdi
	xorl	%eax, %eax
	call	printfPQExpBuffer@PLT
.L1551:
	movq	%rbx, %rdi
	call	PQconninfoFree@PLT
	movq	%r12, %rdi
	call	free@PLT
.L1549:
	xorl	%ebx, %ebx
.L1320:
	movq	120(%rsp), %rdi
	xorq	%fs:40, %rdi
	movq	%rbx, %rax
	jne	.L1567
	addq	$136, %rsp
	.cfi_remember_state
	.cfi_def_cfa_offset 56
	popq	%rbx
	.cfi_def_cfa_offset 48
	popq	%rbp
	.cfi_def_cfa_offset 40
	popq	%r12
	.cfi_def_cfa_offset 32
	popq	%r13
	.cfi_def_cfa_offset 24
	popq	%r14
	.cfi_def_cfa_offset 16
	popq	%r15
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L1392:
	.cfi_restore_state
	addq	$1, %r13
	movzbl	0(%r13), %r14d
	testb	%r14b, %r14b
	jne	.L1421
.L1560:
	movq	%r12, %rdi
	call	free@PLT
.L1556:
	cmpb	$0, 8(%rsp)
	je	.L1320
	movq	%rbp, %rsi
	movq	%rbx, %rdi
	call	conninfo_add_defaults
	testb	%al, %al
	jne	.L1320
	jmp	.L1550
	.p2align 4,,10
	.p2align 3
.L1403:
	movq	%rdx, %rdi
	movq	%rdx, %r14
.L1410:
	movb	$0, (%rdi)
.L1414:
	xorl	%r9d, %r9d
	xorl	%r8d, %r8d
	movq	%rbp, %rcx
	movq	%r13, %rsi
	movq	%rbx, %rdi
	call	conninfo_storeval
	testq	%rax, %rax
	je	.L1551
	movq	%r14, %r13
	jmp	.L1391
	.p2align 4,,10
	.p2align 3
.L1565:
	movzbl	1(%rcx), %eax
	testb	%al, %al
	je	.L1418
	movb	%al, (%rsi)
	addq	$2, %rcx
	movzbl	(%rcx), %eax
	addq	$1, %rsi
	jmp	.L1415
	.p2align 4,,10
	.p2align 3
.L1561:
	movzbl	1(%rsi), %eax
	addq	$2, %rdx
	movb	$0, (%rsi)
	testb	%al, %al
	jne	.L1399
	jmp	.L1400
	.p2align 4,,10
	.p2align 3
.L1568:
	addq	$1, %rdx
	movzbl	(%rdx), %eax
	testb	%al, %al
	je	.L1400
.L1399:
	movzbl	%al, %esi
	testb	$32, 1(%rcx,%rsi,2)
	jne	.L1568
	movq	%rdx, %rsi
	jmp	.L1397
	.p2align 4,,10
	.p2align 3
.L1558:
	call	conninfo_init
	testq	%rax, %rax
	movq	%rax, %rbx
	je	.L1549
	leaq	64(%rsp), %rdi
	movq	%rdi, 24(%rsp)
	call	initPQExpBuffer@PLT
	leaq	96(%rsp), %rax
	movq	%rax, %rdi
	movq	%rax, 16(%rsp)
	call	initPQExpBuffer@PLT
	cmpq	$0, 80(%rsp)
	je	.L1324
	cmpq	$0, 112(%rsp)
	jne	.L1325
.L1324:
	leaq	.LC3(%rip), %rsi
	movq	%rbp, %rdi
	xorl	%eax, %eax
	call	printfPQExpBuffer@PLT
	movq	24(%rsp), %rdi
	call	termPQExpBuffer@PLT
	movq	16(%rsp), %rdi
	call	termPQExpBuffer@PLT
.L1550:
	movq	%rbx, %rdi
	xorl	%ebx, %ebx
	call	PQconninfoFree@PLT
	jmp	.L1320
	.p2align 4,,10
	.p2align 3
.L1400:
	leaq	.LC107(%rip), %rsi
	movq	%r13, %rdx
	movq	%rbp, %rdi
	xorl	%eax, %eax
	call	printfPQExpBuffer@PLT
	jmp	.L1551
	.p2align 4,,10
	.p2align 3
.L1566:
	movb	$0, (%rsi)
	jmp	.L1414
	.p2align 4,,10
	.p2align 3
.L1564:
	movzbl	(%rdx), %eax
	movq	%rdx, %rdi
	movq	%rdx, %rsi
	.p2align 4,,10
	.p2align 3
.L1408:
	testb	%al, %al
	je	.L1569
.L1413:
	movzbl	%al, %r8d
	leaq	1(%rsi), %r14
	testb	$32, 1(%rcx,%r8,2)
	jne	.L1570
	cmpb	$92, %al
	je	.L1571
	movb	%al, (%rdi)
	movzbl	1(%rsi), %eax
	addq	$1, %rdi
	movq	%r14, %rsi
	testb	%al, %al
	jne	.L1413
.L1569:
	movq	%rsi, %r14
	jmp	.L1410
	.p2align 4,,10
	.p2align 3
.L1571:
	movzbl	1(%rsi), %eax
	testb	%al, %al
	je	.L1410
	movb	%al, (%rdi)
	addq	$2, %rsi
	movzbl	(%rsi), %eax
	addq	$1, %rdi
	jmp	.L1408
	.p2align 4,,10
	.p2align 3
.L1570:
	movb	$0, (%rsi)
	jmp	.L1410
	.p2align 4,,10
	.p2align 3
.L1325:
	movq	%r13, %rdi
	call	strdup@PLT
	testq	%rax, %rax
	movq	%rax, 32(%rsp)
	je	.L1324
	movq	%r13, %rdi
	call	uri_prefix_length
	testl	%eax, %eax
	je	.L1572
	movq	32(%rsp), %rdi
	cltq
	leaq	(%rdi,%rax), %rdx
	movzbl	(%rdx), %r15d
	testb	$-65, %r15b
	je	.L1433
	cmpb	$47, %r15b
	je	.L1329
	movq	%rdx, %rax
	jmp	.L1330
	.p2align 4,,10
	.p2align 3
.L1573:
	cmpb	$47, %cl
	je	.L1329
.L1330:
	addq	$1, %rax
	movzbl	(%rax), %ecx
	testb	$-65, %cl
	jne	.L1573
.L1328:
	cmpb	$64, %cl
	jne	.L1329
	cmpb	$64, %r15b
	movq	%rdx, %r10
	je	.L1331
	cmpb	$58, %r15b
	jne	.L1332
	jmp	.L1331
	.p2align 4,,10
	.p2align 3
.L1574:
	cmpb	$64, %r15b
	je	.L1331
.L1332:
	addq	$1, %r10
	movzbl	(%r10), %r15d
	cmpb	$58, %r15b
	jne	.L1574
.L1331:
	movb	$0, (%r10)
	cmpb	$0, (%rdx)
	jne	.L1575
	cmpb	$58, %r15b
	leaq	1(%r10), %rdx
	je	.L1425
.L1336:
	movzbl	1(%r10), %r15d
.L1329:
	movq	%rbx, 40(%rsp)
	movl	%r15d, %ebx
	movq	%r13, 56(%rsp)
	cmpb	$91, %bl
	movabsq	$-8934983331028664319, %r12
	movq	%rdx, %r13
	movq	%rbp, 48(%rsp)
	movq	24(%rsp), %r14
	movq	16(%rsp), %r15
	je	.L1341
	.p2align 4,,10
	.p2align 3
.L1577:
	cmpb	$63, %bl
	movq	%r13, %rbp
	jbe	.L1548
	.p2align 4,,10
	.p2align 3
.L1350:
	addq	$1, %rbp
	movzbl	0(%rbp), %ebx
	cmpb	$63, %bl
	ja	.L1350
.L1548:
	btq	%rbx, %r12
	jnc	.L1350
.L1343:
	movb	$0, 0(%rbp)
	movq	%r13, %rsi
	movq	%r14, %rdi
	call	appendPQExpBufferStr@PLT
	cmpb	$58, %bl
	je	.L1576
.L1352:
	cmpb	$44, %bl
	jne	.L1543
.L1358:
	movl	$44, %esi
	movq	%r14, %rdi
	leaq	1(%rbp), %r13
	call	appendPQExpBufferChar@PLT
	movl	$44, %esi
	movq	%r15, %rdi
	call	appendPQExpBufferChar@PLT
	movzbl	1(%rbp), %ebx
	cmpb	$91, %bl
	jne	.L1577
.L1341:
	movzbl	1(%r13), %eax
	leaq	1(%r13), %rsi
	cmpb	$93, %al
	je	.L1344
	testb	%al, %al
	je	.L1344
	movq	%rsi, %rcx
	jmp	.L1345
	.p2align 4,,10
	.p2align 3
.L1578:
	cmpb	$93, %al
	je	.L1448
	movq	%rdx, %rcx
.L1345:
	movzbl	1(%rcx), %eax
	leaq	1(%rcx), %rdx
	testb	%al, %al
	jne	.L1578
.L1448:
	testb	%al, %al
	je	.L1579
	cmpq	%rdx, %rsi
	je	.L1580
	movzbl	1(%rdx), %ebx
	leaq	2(%rcx), %rbp
	movb	$0, (%rdx)
	cmpb	$63, %bl
	jbe	.L1581
.L1349:
	movq	%rbp, %r14
	movq	56(%rsp), %r13
	movq	48(%rsp), %rbp
	movq	%r14, %rcx
	movq	32(%rsp), %r14
	leaq	.LC100(%rip), %rsi
	movsbl	%bl, %edx
	xorl	%eax, %eax
	movq	40(%rsp), %rbx
	movq	%r13, %r8
	movq	%rbp, %rdi
	subq	%r14, %rcx
	addl	$1, %ecx
	call	printfPQExpBuffer@PLT
	movq	24(%rsp), %rdi
	call	termPQExpBuffer@PLT
	movq	16(%rsp), %rdi
	call	termPQExpBuffer@PLT
	movq	%r14, %rdi
	call	free@PLT
	jmp	.L1550
	.p2align 4,,10
	.p2align 3
.L1443:
	movq	%r13, %rsi
	jmp	.L1393
.L1433:
	movl	%r15d, %ecx
	jmp	.L1328
	.p2align 4,,10
	.p2align 3
.L1581:
	btq	%rbx, %r12
	jnc	.L1349
	movq	%rsi, %r13
	jmp	.L1343
	.p2align 4,,10
	.p2align 3
.L1576:
	movzbl	1(%rbp), %ebx
	leaq	1(%rbp), %r13
	testb	%bl, %bl
	je	.L1440
	cmpb	$47, %bl
	je	.L1440
	cmpb	$63, %bl
	je	.L1354
	cmpb	$44, %bl
	je	.L1355
	movq	%r13, %rbp
	jmp	.L1356
	.p2align 4,,10
	.p2align 3
.L1582:
	cmpb	$44, %bl
	je	.L1353
.L1356:
	addq	$1, %rbp
	movzbl	0(%rbp), %ebx
	testb	%bl, %bl
	je	.L1353
	cmpb	$47, %bl
	je	.L1353
	cmpb	$63, %bl
	jne	.L1582
.L1353:
	movb	$0, 0(%rbp)
	movq	%r13, %rsi
	movq	%r15, %rdi
	call	appendPQExpBufferStr@PLT
	jmp	.L1352
.L1344:
	testb	%al, %al
	movq	40(%rsp), %rbx
	movq	56(%rsp), %r13
	movq	48(%rsp), %rbp
	jne	.L1427
.L1426:
	leaq	.LC98(%rip), %rsi
	movq	%r13, %rdx
	movq	%rbp, %rdi
	xorl	%eax, %eax
	call	printfPQExpBuffer@PLT
.L1383:
	movq	24(%rsp), %rdi
	call	termPQExpBuffer@PLT
	movq	16(%rsp), %rdi
	call	termPQExpBuffer@PLT
	movq	32(%rsp), %rdi
	call	free@PLT
	jmp	.L1550
.L1572:
	leaq	.LC96(%rip), %rsi
	movq	%r13, %rdx
	movq	%rbp, %rdi
	call	printfPQExpBuffer@PLT
	jmp	.L1383
.L1543:
	movq	%rbp, %r14
	movl	%ebx, %r15d
	movq	48(%rsp), %rbp
	movq	40(%rsp), %rbx
.L1357:
	cmpq	$0, 80(%rsp)
	je	.L1383
	cmpq	$0, 112(%rsp)
	je	.L1383
	movq	64(%rsp), %rdx
	cmpb	$0, (%rdx)
	jne	.L1362
.L1366:
	movq	96(%rsp), %rdx
	cmpb	$0, (%rdx)
	jne	.L1583
.L1364:
	testb	%r15b, %r15b
	je	.L1449
	cmpb	$63, %r15b
	je	.L1449
	movzbl	1(%r14), %r15d
	leaq	1(%r14), %rdx
	movq	%rdx, %r13
	testb	%r15b, %r15b
	jne	.L1557
	jmp	.L1371
	.p2align 4,,10
	.p2align 3
.L1584:
	addq	$1, %r13
	movzbl	0(%r13), %r15d
	testb	%r15b, %r15b
	je	.L1371
.L1557:
	cmpb	$63, %r15b
	jne	.L1584
.L1371:
	movb	$0, 0(%r13)
	cmpb	$0, 1(%r14)
	jne	.L1585
.L1369:
	testb	%r15b, %r15b
	jne	.L1586
.L1388:
	movq	24(%rsp), %rdi
	call	termPQExpBuffer@PLT
	movq	16(%rsp), %rdi
	call	termPQExpBuffer@PLT
	movq	32(%rsp), %rdi
	call	free@PLT
	jmp	.L1556
.L1559:
	leaq	.LC3(%rip), %rsi
	movq	%rbp, %rdi
	xorl	%eax, %eax
	call	printfPQExpBuffer@PLT
	movq	%rbx, %rdi
	xorl	%ebx, %ebx
	call	PQconninfoFree@PLT
	jmp	.L1320
.L1580:
	movq	40(%rsp), %rbx
	movq	56(%rsp), %r13
	movq	48(%rsp), %rbp
.L1427:
	leaq	.LC99(%rip), %rsi
	movq	%r13, %rdx
	movq	%rbp, %rdi
	xorl	%eax, %eax
	call	printfPQExpBuffer@PLT
	jmp	.L1383
.L1440:
	movq	%r13, %rbp
	jmp	.L1353
.L1354:
	movq	16(%rsp), %rdi
	movq	%rbp, %r14
	movq	%r13, %rsi
	movb	$0, 1(%r14)
	movl	%ebx, %r15d
	movq	48(%rsp), %rbp
	movq	40(%rsp), %rbx
	movq	%r13, 40(%rsp)
	call	appendPQExpBufferStr@PLT
	movq	40(%rsp), %rax
	movq	%rax, %r14
	jmp	.L1357
.L1355:
	movb	$0, 1(%rbp)
	movq	%r13, %rsi
	movq	%r15, %rdi
	call	appendPQExpBufferStr@PLT
	movq	%r13, %rbp
	jmp	.L1358
.L1575:
	leaq	.LC95(%rip), %rsi
	xorl	%r8d, %r8d
	movl	$1, %r9d
	movq	%rbp, %rcx
	movq	%rbx, %rdi
	movq	%r10, 40(%rsp)
	call	conninfo_storeval
	testq	%rax, %rax
	je	.L1383
	movq	40(%rsp), %r10
	cmpb	$58, %r15b
	leaq	1(%r10), %rdx
	jne	.L1336
	cmpb	$64, (%r10)
	movq	%r10, %r12
	jne	.L1425
.L1337:
	movb	$0, (%r12)
	cmpb	$0, 1(%r10)
	jne	.L1339
.L1547:
	leaq	1(%r12), %rdx
	movq	%r12, %r10
	jmp	.L1336
.L1425:
	movq	%r10, %r12
.L1338:
	addq	$1, %r12
	cmpb	$64, (%r12)
	jne	.L1338
	jmp	.L1337
.L1449:
	movq	%r14, %r13
	jmp	.L1369
.L1583:
	leaq	.LC33(%rip), %rsi
	xorl	%r8d, %r8d
	movl	$1, %r9d
	movq	%rbp, %rcx
	movq	%rbx, %rdi
	call	conninfo_storeval
	testq	%rax, %rax
	jne	.L1364
	jmp	.L1383
.L1362:
	leaq	.LC101(%rip), %rsi
	xorl	%r8d, %r8d
	movl	$1, %r9d
	movq	%rbp, %rcx
	movq	%rbx, %rdi
	call	conninfo_storeval
	testq	%rax, %rax
	jne	.L1366
	jmp	.L1383
.L1586:
	leaq	1(%r13), %rdx
	leaq	.LC14(%rip), %r15
	movq	%rdx, %rdi
.L1376:
	movzbl	(%rdi), %ecx
	testb	%cl, %cl
	je	.L1388
	movq	%rdi, %rax
	xorl	%r14d, %r14d
	jmp	.L1389
.L1377:
	cmpb	$38, %cl
	je	.L1450
	testb	%cl, %cl
	je	.L1450
	leaq	1(%rax), %rdx
.L1379:
	movzbl	1(%rax), %ecx
	movq	%rdx, %rax
.L1389:
	cmpb	$61, %cl
	movq	%rax, %r12
	jne	.L1377
	testq	%r14, %r14
	jne	.L1587
	leaq	1(%rax), %rdx
	movb	$0, (%rax)
	movq	%rdx, %r14
	jmp	.L1379
.L1450:
	testb	%cl, %cl
	je	.L1381
	leaq	1(%rax), %r12
	movb	$0, (%rax)
.L1381:
	testq	%r14, %r14
	je	.L1588
	movq	%rbp, %rsi
	call	conninfo_uri_decode
	testq	%rax, %rax
	movq	%rax, %r13
	je	.L1383
	movq	%r14, %rdi
	movq	%rbp, %rsi
	call	conninfo_uri_decode
	testq	%rax, %rax
	movq	%rax, %r14
	je	.L1589
	leaq	.LC105(%rip), %rdi
	movl	$4, %ecx
	movq	%r13, %rsi
	repz cmpsb
	seta	%al
	sbbb	$0, %al
	testb	%al, %al
	jne	.L1385
	leaq	.LC106(%rip), %rdi
	movl	$5, %ecx
	movq	%r14, %rsi
	repz cmpsb
	seta	%al
	sbbb	$0, %al
	testb	%al, %al
	je	.L1590
.L1385:
	xorl	%r9d, %r9d
	movl	$1, %r8d
	movq	%rbp, %rcx
	movq	%r14, %rdx
	movq	%r13, %rsi
	movq	%rbx, %rdi
	call	conninfo_storeval
	testq	%rax, %rax
	jne	.L1591
	cmpq	$0, 8(%rbp)
	movq	%r13, %r10
	movq	%r14, %r15
	je	.L1592
.L1429:
	movq	%r10, %rdi
	call	free@PLT
	movq	%r15, %rdi
	call	free@PLT
	jmp	.L1383
	.p2align 4,,10
	.p2align 3
.L1339:
	leaq	.LC97(%rip), %rsi
	xorl	%r8d, %r8d
	movl	$1, %r9d
	movq	%rbp, %rcx
	movq	%rbx, %rdi
	call	conninfo_storeval
	testq	%rax, %rax
	jne	.L1547
	jmp	.L1383
.L1587:
	leaq	.LC103(%rip), %rsi
	movq	%rdi, %rdx
	xorl	%eax, %eax
	movq	%rbp, %rdi
	call	printfPQExpBuffer@PLT
	jmp	.L1383
.L1585:
	leaq	.LC102(%rip), %rsi
	xorl	%r8d, %r8d
	movl	$1, %r9d
	movq	%rbp, %rcx
	movq	%rbx, %rdi
	call	conninfo_storeval
	testq	%rax, %rax
	jne	.L1369
	jmp	.L1383
	.p2align 4,,10
	.p2align 3
.L1591:
	movq	%r13, %rdi
	call	free@PLT
	movq	%r14, %rdi
	call	free@PLT
.L1387:
	movq	%r12, %rdi
	jmp	.L1376
.L1592:
	leaq	.LC109(%rip), %rsi
	movq	%r13, %rdx
	movq	%rbp, %rdi
	xorl	%eax, %eax
	movq	%r13, 8(%rsp)
	call	printfPQExpBuffer@PLT
	movq	8(%rsp), %r10
	jmp	.L1429
.L1589:
	movq	%r13, %rdi
	call	free@PLT
	jmp	.L1383
.L1588:
	leaq	.LC104(%rip), %rsi
	movq	%rdi, %rdx
	xorl	%eax, %eax
	movq	%rbp, %rdi
	call	printfPQExpBuffer@PLT
	jmp	.L1383
.L1579:
	movq	40(%rsp), %rbx
	movq	56(%rsp), %r13
	movq	48(%rsp), %rbp
	jmp	.L1426
.L1590:
	movq	%r13, %rdi
	call	free@PLT
	movq	%r14, %rdi
	call	free@PLT
	leaq	.LC15(%rip), %rsi
	xorl	%r9d, %r9d
	movl	$1, %r8d
	movq	%rbp, %rcx
	movq	%r15, %rdx
	movq	%rbx, %rdi
	call	conninfo_storeval
	testq	%rax, %rax
	jne	.L1387
	cmpq	$0, 8(%rbp)
	jne	.L1383
	leaq	.LC15(%rip), %rdx
	leaq	.LC109(%rip), %rsi
	movq	%rbp, %rdi
	xorl	%eax, %eax
	call	printfPQExpBuffer@PLT
	jmp	.L1383
.L1567:
	call	__stack_chk_fail@PLT
	.cfi_endproc
.LFE958:
	.size	parse_connection_string, .-parse_connection_string
	.p2align 4,,15
	.type	connectOptions1, @function
connectOptions1:
.LFB918:
	.cfi_startproc
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset 6, -16
	pushq	%rbx
	.cfi_def_cfa_offset 24
	.cfi_offset 3, -24
	movq	%rdi, %rbx
	movq	%rsi, %rdi
	leaq	928(%rbx), %rsi
	movl	$1, %edx
	subq	$24, %rsp
	.cfi_def_cfa_offset 48
	call	parse_connection_string
	testq	%rax, %rax
	je	.L1598
	movq	%rax, %rsi
	movq	%rbx, %rdi
	movq	%rax, %rbp
	call	fillPGconn
	testb	%al, %al
	movb	%al, 15(%rsp)
	je	.L1599
	movq	%rbp, %rdi
	call	PQconninfoFree@PLT
	movzbl	15(%rsp), %eax
	addq	$24, %rsp
	.cfi_remember_state
	.cfi_def_cfa_offset 24
	popq	%rbx
	.cfi_def_cfa_offset 16
	popq	%rbp
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L1599:
	.cfi_restore_state
	movl	$1, 336(%rbx)
	movq	%rbp, %rdi
	call	PQconninfoFree@PLT
	movzbl	15(%rsp), %eax
	addq	$24, %rsp
	.cfi_remember_state
	.cfi_def_cfa_offset 24
	popq	%rbx
	.cfi_def_cfa_offset 16
	popq	%rbp
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L1598:
	.cfi_restore_state
	movl	$1, 336(%rbx)
	addq	$24, %rsp
	.cfi_def_cfa_offset 24
	xorl	%eax, %eax
	popq	%rbx
	.cfi_def_cfa_offset 16
	popq	%rbp
	.cfi_def_cfa_offset 8
	ret
	.cfi_endproc
.LFE918:
	.size	connectOptions1, .-connectOptions1
	.p2align 4,,15
	.globl	PQconninfoParse
	.type	PQconninfoParse, @function
PQconninfoParse:
.LFB956:
	.cfi_startproc
	pushq	%r12
	.cfi_def_cfa_offset 16
	.cfi_offset 12, -16
	pushq	%rbp
	.cfi_def_cfa_offset 24
	.cfi_offset 6, -24
	movq	%rdi, %r12
	pushq	%rbx
	.cfi_def_cfa_offset 32
	.cfi_offset 3, -32
	subq	$48, %rsp
	.cfi_def_cfa_offset 80
	movq	%fs:40, %rax
	movq	%rax, 40(%rsp)
	xorl	%eax, %eax
	testq	%rsi, %rsi
	je	.L1601
	leaq	16(%rsp), %rbp
	movq	$0, (%rsi)
	movq	%rsi, %rbx
	movq	%rbp, %rdi
	call	initPQExpBuffer@PLT
	cmpq	$0, 32(%rsp)
	je	.L1605
	xorl	%edx, %edx
	movq	%rbp, %rsi
	movq	%r12, %rdi
	call	parse_connection_string
	testq	%rax, %rax
	jne	.L1604
	movq	16(%rsp), %rdx
	movq	%rdx, (%rbx)
.L1600:
	movq	40(%rsp), %rcx
	xorq	%fs:40, %rcx
	jne	.L1610
	addq	$48, %rsp
	.cfi_remember_state
	.cfi_def_cfa_offset 32
	popq	%rbx
	.cfi_def_cfa_offset 24
	popq	%rbp
	.cfi_def_cfa_offset 16
	popq	%r12
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L1601:
	.cfi_restore_state
	leaq	16(%rsp), %rbp
	movq	%rbp, %rdi
	call	initPQExpBuffer@PLT
	cmpq	$0, 32(%rsp)
	je	.L1605
	xorl	%edx, %edx
	movq	%rbp, %rsi
	movq	%r12, %rdi
	call	parse_connection_string
.L1604:
	movq	%rbp, %rdi
	movq	%rax, 8(%rsp)
	call	termPQExpBuffer@PLT
	movq	8(%rsp), %rax
	jmp	.L1600
	.p2align 4,,10
	.p2align 3
.L1605:
	xorl	%eax, %eax
	jmp	.L1600
.L1610:
	call	__stack_chk_fail@PLT
	.cfi_endproc
.LFE956:
	.size	PQconninfoParse, .-PQconninfoParse
	.p2align 4,,15
	.globl	PQconndefaults
	.type	PQconndefaults, @function
PQconndefaults:
.LFB922:
	.cfi_startproc
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset 6, -16
	pushq	%rbx
	.cfi_def_cfa_offset 24
	.cfi_offset 3, -24
	subq	$40, %rsp
	.cfi_def_cfa_offset 64
	movq	%rsp, %rbx
	movq	%rbx, %rdi
	movq	%fs:40, %rax
	movq	%rax, 24(%rsp)
	xorl	%eax, %eax
	call	initPQExpBuffer@PLT
	cmpq	$0, 16(%rsp)
	je	.L1615
	movq	%rbx, %rdi
	call	conninfo_init
	testq	%rax, %rax
	movq	%rax, %rbp
	je	.L1613
	xorl	%esi, %esi
	movq	%rax, %rdi
	call	conninfo_add_defaults
	testb	%al, %al
	je	.L1620
.L1613:
	movq	%rbx, %rdi
	call	termPQExpBuffer@PLT
.L1611:
	movq	24(%rsp), %rdx
	xorq	%fs:40, %rdx
	movq	%rbp, %rax
	jne	.L1621
	addq	$40, %rsp
	.cfi_remember_state
	.cfi_def_cfa_offset 24
	popq	%rbx
	.cfi_def_cfa_offset 16
	popq	%rbp
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L1620:
	.cfi_restore_state
	movq	%rbp, %rdi
	xorl	%ebp, %ebp
	call	PQconninfoFree@PLT
	jmp	.L1613
	.p2align 4,,10
	.p2align 3
.L1615:
	xorl	%ebp, %ebp
	jmp	.L1611
.L1621:
	call	__stack_chk_fail@PLT
	.cfi_endproc
.LFE922:
	.size	PQconndefaults, .-PQconndefaults
	.section	.rodata.str1.1
.LC110:
	.string	"localhost"
	.section	.rodata.str1.8
	.align 8
.LC111:
	.string	"could not match %d host names to %d hostaddr values\n"
	.section	.rodata.str1.1
.LC112:
	.string	"/tmp"
	.section	.rodata.str1.8
	.align 8
.LC113:
	.string	"could not match %d port numbers to %d hosts\n"
	.section	.rodata.str1.1
.LC114:
	.string	".pgpass"
	.section	.rodata.str1.8
	.align 8
.LC115:
	.string	"WARNING: password file \"%s\" is not a plain file\n"
	.align 8
.LC116:
	.string	"WARNING: password file \"%s\" has group or world access; permissions should be u=rw (0600) or less\n"
	.align 8
.LC117:
	.string	"WARNING: line %d too long in password file \"%s\"\n"
	.section	.rodata.str1.1
.LC118:
	.string	"disable"
	.section	.rodata.str1.8
	.align 8
.LC119:
	.string	"invalid channel_binding value: \"%s\"\n"
	.section	.rodata.str1.1
.LC120:
	.string	"allow"
.LC121:
	.string	"verify-ca"
.LC122:
	.string	"verify-full"
.LC123:
	.string	"invalid sslmode value: \"%s\"\n"
	.section	.rodata.str1.8
	.align 8
.LC124:
	.string	"invalid ssl_min_protocol_version value: \"%s\"\n"
	.align 8
.LC125:
	.string	"invalid ssl_max_protocol_version value: \"%s\"\n"
	.align 8
.LC126:
	.string	"invalid SSL protocol version range\n"
	.align 8
.LC127:
	.string	"invalid gssencmode value: \"%s\"\n"
	.align 8
.LC128:
	.string	"gssencmode value \"%s\" invalid when GSSAPI support is not compiled in\n"
	.section	.rodata.str1.1
.LC129:
	.string	"any"
	.section	.rodata.str1.8
	.align 8
.LC130:
	.string	"invalid target_session_attrs value: \"%s\"\n"
	.text
	.p2align 4,,15
	.type	connectOptions2, @function
connectOptions2:
.LFB921:
	.cfi_startproc
	pushq	%r15
	.cfi_def_cfa_offset 16
	.cfi_offset 15, -16
	pushq	%r14
	.cfi_def_cfa_offset 24
	.cfi_offset 14, -24
	movq	%rdi, %r15
	pushq	%r13
	.cfi_def_cfa_offset 32
	.cfi_offset 13, -32
	pushq	%r12
	.cfi_def_cfa_offset 40
	.cfi_offset 12, -40
	pushq	%rbp
	.cfi_def_cfa_offset 48
	.cfi_offset 6, -48
	pushq	%rbx
	.cfi_def_cfa_offset 56
	.cfi_offset 3, -56
	subq	$1576, %rsp
	.cfi_def_cfa_offset 1632
	movq	16(%rdi), %r12
	movl	$0, 396(%rdi)
	movq	%fs:40, %rax
	movq	%rax, 1560(%rsp)
	xorl	%eax, %eax
	testq	%r12, %r12
	je	.L1623
	movzbl	(%r12), %eax
	testb	%al, %al
	je	.L1623
	movq	%r12, %rdx
	movl	$1, %ebx
	.p2align 4,,10
	.p2align 3
.L1625:
	cmpb	$44, %al
	sete	%al
	addq	$1, %rdx
	movzbl	%al, %eax
	addl	%eax, %ebx
	movzbl	(%rdx), %eax
	testb	%al, %al
	jne	.L1625
	movl	%ebx, 392(%r15)
	movslq	%ebx, %rdi
	movl	$40, %esi
	call	calloc@PLT
	testq	%rax, %rax
	movq	%rax, 400(%r15)
	je	.L1637
.L1627:
	cmpb	$0, (%r12)
	je	.L1634
	testl	%ebx, %ebx
	movq	%r12, 56(%rsp)
	movb	$1, 55(%rsp)
	jle	.L1634
	movq	400(%r15), %rax
	xorl	%r14d, %r14d
	leaq	55(%rsp), %r12
	leaq	56(%rsp), %rbp
	jmp	.L1636
	.p2align 4,,10
	.p2align 3
.L1948:
	movl	392(%r15), %ebx
	addl	$1, %r14d
	cmpl	%r14d, %ebx
	jle	.L1634
	cmpb	$0, 55(%rsp)
	je	.L1634
.L1636:
	movslq	%r14d, %rdx
	movq	%r12, %rsi
	movq	%rbp, %rdi
	leaq	(%rdx,%rdx,4), %r13
	salq	$3, %r13
	leaq	(%rax,%r13), %rbx
	call	parse_comma_separated_list
	movq	%rax, 16(%rbx)
	movq	400(%r15), %rax
	cmpq	$0, 16(%rax,%r13)
	jne	.L1948
	.p2align 4,,10
	.p2align 3
.L1637:
	leaq	928(%r15), %rdi
	leaq	.LC3(%rip), %rsi
	movl	$1, 336(%r15)
	xorl	%eax, %eax
	xorl	%ebx, %ebx
	call	printfPQExpBuffer@PLT
.L1622:
	movq	1560(%rsp), %rsi
	xorq	%fs:40, %rsi
	movl	%ebx, %eax
	jne	.L1949
	addq	$1576, %rsp
	.cfi_remember_state
	.cfi_def_cfa_offset 56
	popq	%rbx
	.cfi_def_cfa_offset 48
	popq	%rbp
	.cfi_def_cfa_offset 40
	popq	%r12
	.cfi_def_cfa_offset 32
	popq	%r13
	.cfi_def_cfa_offset 24
	popq	%r14
	.cfi_def_cfa_offset 16
	popq	%r15
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L1623:
	.cfi_restore_state
	movq	8(%r15), %rbp
	testq	%rbp, %rbp
	je	.L1628
	movzbl	0(%rbp), %eax
	testb	%al, %al
	je	.L1628
	movq	%rbp, %rdx
	movl	$1, %ebx
	.p2align 4,,10
	.p2align 3
.L1630:
	cmpb	$44, %al
	sete	%al
	addq	$1, %rdx
	movzbl	%al, %eax
	addl	%eax, %ebx
	movzbl	(%rdx), %eax
	testb	%al, %al
	jne	.L1630
	movl	%ebx, 392(%r15)
	movslq	%ebx, %rdi
	jmp	.L1631
	.p2align 4,,10
	.p2align 3
.L1628:
	movl	$1, 392(%r15)
	movl	$1, %edi
	movl	$1, %ebx
.L1631:
	movl	$40, %esi
	call	calloc@PLT
	testq	%rax, %rax
	movq	%rax, 400(%r15)
	je	.L1637
	testq	%r12, %r12
	jne	.L1627
	testq	%rbp, %rbp
	je	.L1638
.L1951:
	cmpb	$0, 0(%rbp)
	je	.L1638
	testl	%ebx, %ebx
	movq	%rbp, 56(%rsp)
	movb	$1, 55(%rsp)
	jle	.L1741
	movq	400(%r15), %rax
	xorl	%r14d, %r14d
	leaq	55(%rsp), %r12
	leaq	56(%rsp), %rbp
	jmp	.L1642
	.p2align 4,,10
	.p2align 3
.L1950:
	movl	392(%r15), %ebx
	addl	$1, %r14d
	movzbl	55(%rsp), %ecx
	cmpl	%r14d, %ebx
	jle	.L1641
	testb	%cl, %cl
	je	.L1754
.L1642:
	movslq	%r14d, %rdx
	movq	%r12, %rsi
	movq	%rbp, %rdi
	leaq	(%rdx,%rdx,4), %r13
	salq	$3, %r13
	leaq	(%rax,%r13), %rbx
	call	parse_comma_separated_list
	movq	%rax, 8(%rbx)
	movq	400(%r15), %rax
	cmpq	$0, 8(%rax,%r13)
	movq	%rax, %rdx
	jne	.L1950
	jmp	.L1637
	.p2align 4,,10
	.p2align 3
.L1634:
	movq	8(%r15), %rbp
	testq	%rbp, %rbp
	jne	.L1951
.L1638:
	testl	%ebx, %ebx
	jle	.L1650
	movq	400(%r15), %rdx
.L1649:
	xorl	%ebp, %ebp
	leaq	.LC112(%rip), %r12
	jmp	.L1657
	.p2align 4,,10
	.p2align 3
.L1952:
	movl	$1, (%rbx)
.L1652:
	addl	$1, %ebp
	cmpl	%ebp, 392(%r15)
	jle	.L1650
.L1953:
	movq	400(%r15), %rdx
.L1657:
	movslq	%ebp, %rax
	leaq	(%rax,%rax,4), %rax
	leaq	(%rdx,%rax,8), %rbx
	movq	16(%rbx), %rax
	testq	%rax, %rax
	je	.L1651
	cmpb	$0, (%rax)
	jne	.L1952
.L1651:
	movq	8(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L1653
	cmpb	$0, (%rdi)
	je	.L1654
	movl	$0, (%rbx)
	cmpb	$47, (%rdi)
	jne	.L1652
	movl	$2, (%rbx)
	addl	$1, %ebp
	cmpl	%ebp, 392(%r15)
	jg	.L1953
	.p2align 4,,10
	.p2align 3
.L1650:
	movq	24(%r15), %rax
	testq	%rax, %rax
	je	.L1658
	cmpb	$0, (%rax)
	je	.L1658
	movl	392(%r15), %ecx
	movq	%rax, 56(%rsp)
	movb	$1, 55(%rsp)
	testl	%ecx, %ecx
	jle	.L1744
	movq	400(%r15), %r14
	xorl	%r13d, %r13d
	leaq	55(%rsp), %r12
	leaq	56(%rsp), %rbp
	jmp	.L1662
	.p2align 4,,10
	.p2align 3
.L1955:
	movl	392(%r15), %ecx
	addl	$1, %r13d
	movzbl	55(%rsp), %eax
	cmpl	%r13d, %ecx
	jle	.L1661
	testb	%al, %al
	je	.L1954
.L1662:
	movslq	%r13d, %rax
	movq	%r12, %rsi
	movq	%rbp, %rdi
	leaq	(%rax,%rax,4), %rbx
	call	parse_comma_separated_list
	salq	$3, %rbx
	addq	%rbx, %r14
	movq	%rax, 24(%r14)
	movq	400(%r15), %r14
	cmpq	$0, 24(%r14,%rbx)
	jne	.L1955
	jmp	.L1637
	.p2align 4,,10
	.p2align 3
.L1654:
	call	free@PLT
.L1653:
	movq	%r12, %rdi
	call	strdup@PLT
	testq	%rax, %rax
	movq	%rax, 8(%rbx)
	movl	$2, (%rbx)
	jne	.L1652
	jmp	.L1637
.L1967:
	testb	%al, %al
	jne	.L1743
	subl	$1, %ecx
	jg	.L1667
	.p2align 4,,10
	.p2align 3
.L1658:
	movq	104(%r15), %rdi
	testq	%rdi, %rdi
	je	.L1672
	cmpb	$0, (%rdi)
	jne	.L1673
	call	free@PLT
.L1672:
	leaq	928(%r15), %rdi
	call	pg_fe_getauthname@PLT
	testq	%rax, %rax
	movq	%rax, 104(%r15)
	je	.L1956
.L1673:
	movq	88(%r15), %rdi
	testq	%rdi, %rdi
	je	.L1674
	cmpb	$0, (%rdi)
	jne	.L1678
	call	free@PLT
.L1674:
	movq	104(%r15), %rdi
	call	strdup@PLT
	testq	%rax, %rax
	movq	%rax, 88(%r15)
	je	.L1637
.L1678:
	movq	112(%r15), %rax
	testq	%rax, %rax
	je	.L1676
	cmpb	$0, (%rax)
	jne	.L1679
.L1676:
	movq	120(%r15), %rax
	testq	%rax, %rax
	je	.L1680
	cmpb	$0, (%rax)
	je	.L1680
.L1681:
	movl	392(%r15), %eax
	testl	%eax, %eax
	jle	.L1679
	leaq	528(%rsp), %r14
	movl	$0, 4(%rsp)
	.p2align 4,,10
	.p2align 3
.L1718:
	movslq	4(%rsp), %rax
	leaq	(%rax,%rax,4), %rdx
	movq	400(%r15), %rax
	leaq	(%rax,%rdx,8), %rax
	movq	8(%rax), %r13
	movq	%rax, 8(%rsp)
	testq	%r13, %r13
	je	.L1686
	cmpb	$0, 0(%r13)
	jne	.L1687
.L1686:
	movq	8(%rsp), %rax
	movq	16(%rax), %r13
.L1687:
	movq	104(%r15), %rax
	movq	%rax, 24(%rsp)
	movq	88(%r15), %rax
	testq	%rax, %rax
	movq	%rax, 16(%rsp)
	je	.L1946
	movq	16(%rsp), %rax
	cmpb	$0, (%rax)
	je	.L1946
	movq	24(%rsp), %rax
	testq	%rax, %rax
	je	.L1946
	cmpb	$0, (%rax)
	je	.L1946
	testq	%r13, %r13
	je	.L1749
	movzbl	0(%r13), %eax
	testb	%al, %al
	je	.L1749
	cmpb	$47, %al
	jne	.L1691
	leaq	.LC112(%rip), %rdi
	movq	%r13, %rsi
	movl	$5, %ecx
	repz cmpsb
	seta	%al
	sbbb	$0, %al
	testb	%al, %al
	leaq	.LC110(%rip), %rax
	cmove	%rax, %r13
	.p2align 4,,10
	.p2align 3
.L1691:
	movq	8(%rsp), %rax
	movq	24(%rax), %rax
	testq	%rax, %rax
	je	.L1751
	cmpb	$0, (%rax)
	leaq	.LC26(%rip), %rsi
	cmovne	%rax, %rsi
	movq	%rsi, 40(%rsp)
.L1692:
	movq	120(%r15), %rax
	leaq	64(%rsp), %rdx
	movl	$1, %edi
	movq	%rax, %rsi
	movq	%rax, 32(%rsp)
	call	__xstat@PLT
	testl	%eax, %eax
	movl	%eax, %r12d
	jne	.L1946
	movl	88(%rsp), %eax
	movl	%eax, %edx
	andl	$61440, %edx
	cmpl	$32768, %edx
	jne	.L1957
	testb	$63, %al
	jne	.L1958
	movq	32(%rsp), %rdi
	leaq	.LC4(%rip), %rsi
	call	fopen@PLT
	testq	%rax, %rax
	movq	%rax, %rbx
	je	.L1946
	leaq	208(%rsp), %rbp
	.p2align 4,,10
	.p2align 3
.L1695:
	movq	%rbx, %rdi
	call	feof@PLT
	testl	%eax, %eax
	jne	.L1716
	movq	%rbx, %rdi
	call	ferror@PLT
	testl	%eax, %eax
	jne	.L1716
	movq	%rbx, %rdx
	movl	$320, %esi
	movq	%rbp, %rdi
	call	fgets@PLT
	testq	%rax, %rax
	je	.L1716
	addl	$1, %r12d
	movq	%rbp, %rax
.L1697:
	movl	(%rax), %ecx
	addq	$4, %rax
	leal	-16843009(%rcx), %edx
	notl	%ecx
	andl	%ecx, %edx
	andl	$-2139062144, %edx
	je	.L1697
	movl	%edx, %ecx
	shrl	$16, %ecx
	testl	$32896, %edx
	cmove	%ecx, %edx
	leaq	2(%rax), %rcx
	movl	%edx, %esi
	cmove	%rcx, %rax
	addb	%dl, %sil
	movzbl	208(%rsp), %edx
	sbbq	$3, %rax
	subq	%rbp, %rax
	cmpl	$318, %eax
	jbe	.L1699
	subl	$1, %eax
	cltq
	cmpb	$10, 208(%rsp,%rax)
	je	.L1699
	cmpb	$35, %dl
	je	.L1705
	movq	stderr@GOTPCREL(%rip), %rax
	movq	32(%rsp), %rcx
	leaq	.LC117(%rip), %rsi
	movl	%r12d, %edx
	movq	(%rax), %rdi
	xorl	%eax, %eax
	call	pg_fprintf@PLT
	jmp	.L1705
	.p2align 4,,10
	.p2align 3
.L1959:
	movq	%rbx, %rdi
	call	ferror@PLT
	testl	%eax, %eax
	jne	.L1704
	movq	%rbx, %rdx
	movl	$320, %esi
	movq	%r14, %rdi
	call	fgets@PLT
	testq	%rax, %rax
	je	.L1704
	movq	%r14, %rax
.L1702:
	movl	(%rax), %ecx
	addq	$4, %rax
	leal	-16843009(%rcx), %edx
	notl	%ecx
	andl	%ecx, %edx
	andl	$-2139062144, %edx
	je	.L1702
	movl	%edx, %ecx
	shrl	$16, %ecx
	testl	$32896, %edx
	cmove	%ecx, %edx
	leaq	2(%rax), %rcx
	movl	%edx, %esi
	cmove	%rcx, %rax
	addb	%dl, %sil
	sbbq	$3, %rax
	subq	%r14, %rax
	cmpl	$318, %eax
	jbe	.L1704
	subl	$1, %eax
	cltq
	cmpb	$10, 528(%rsp,%rax)
	je	.L1704
.L1705:
	movq	%rbx, %rdi
	call	feof@PLT
	testl	%eax, %eax
	je	.L1959
.L1704:
	movzbl	208(%rsp), %edx
.L1699:
	cmpb	$35, %dl
	je	.L1695
	movq	%rbp, %rdi
	call	pg_strip_crlf@PLT
	testl	%eax, %eax
	je	.L1695
	movq	%r13, %rsi
	movq	%rbp, %rdi
	call	pwdfMatchesString
	testq	%rax, %rax
	je	.L1695
	movq	40(%rsp), %rsi
	movq	%rax, %rdi
	call	pwdfMatchesString
	testq	%rax, %rax
	je	.L1695
	movq	16(%rsp), %rsi
	movq	%rax, %rdi
	call	pwdfMatchesString
	testq	%rax, %rax
	je	.L1695
	movq	24(%rsp), %rsi
	movq	%rax, %rdi
	call	pwdfMatchesString
	testq	%rax, %rax
	je	.L1695
	movq	%rax, %rdi
	call	strdup@PLT
	movq	%rbx, %rdi
	movq	%rax, %r12
	call	fclose@PLT
	testq	%r12, %r12
	je	.L1708
	movzbl	(%r12), %eax
	cmpb	$58, %al
	je	.L1755
	testb	%al, %al
	movq	%r12, %rsi
	movq	%r12, %rdx
	jne	.L1709
	jmp	.L1755
	.p2align 4,,10
	.p2align 3
.L1945:
	movq	%rcx, %rdi
	movzbl	(%rdx), %eax
	movq	%rdx, %rcx
	movq	%rdi, %rdx
.L1713:
	movb	%al, (%rsi)
	movzbl	1(%rcx), %eax
	addq	$1, %rsi
	cmpb	$58, %al
	je	.L1711
	testb	%al, %al
	je	.L1711
.L1709:
	cmpb	$92, %al
	leaq	1(%rdx), %rcx
	jne	.L1945
	movzbl	1(%rdx), %eax
	testb	%al, %al
	je	.L1945
	addq	$2, %rdx
	jmp	.L1713
	.p2align 4,,10
	.p2align 3
.L1680:
	leaq	528(%rsp), %rbx
	movl	$1024, %esi
	movq	%rbx, %rdi
	call	pqGetHomeDirectory@PLT
	testb	%al, %al
	je	.L1682
	movq	120(%r15), %rdi
	testq	%rdi, %rdi
	je	.L1683
	call	free@PLT
.L1683:
	movl	$1024, %edi
	call	malloc@PLT
	testq	%rax, %rax
	movq	%rax, 120(%r15)
	je	.L1637
	leaq	.LC114(%rip), %r8
	leaq	.LC90(%rip), %rdx
	movq	%rax, %rdi
	movq	%rbx, %rcx
	movl	$1024, %esi
	xorl	%eax, %eax
	call	pg_snprintf@PLT
.L1682:
	movq	120(%r15), %rax
	testq	%rax, %rax
	jne	.L1960
.L1679:
	movq	128(%r15), %rdx
	testq	%rdx, %rdx
	je	.L1719
	leaq	.LC118(%rip), %rdi
	movl	$8, %ecx
	movq	%rdx, %rsi
	repz cmpsb
	seta	%al
	sbbb	$0, %al
	testb	%al, %al
	je	.L1723
	leaq	.LC16(%rip), %rdi
	movl	$7, %ecx
	movq	%rdx, %rsi
	repz cmpsb
	seta	%al
	sbbb	$0, %al
	testb	%al, %al
	jne	.L1961
.L1723:
	movq	168(%r15), %rdx
	testq	%rdx, %rdx
	je	.L1962
	leaq	.LC118(%rip), %rdi
	movl	$8, %ecx
	movq	%rdx, %rsi
	repz cmpsb
	seta	%al
	sbbb	$0, %al
	testb	%al, %al
	je	.L1727
	leaq	.LC120(%rip), %rdi
	movl	$6, %ecx
	movq	%rdx, %rsi
	repz cmpsb
	seta	%al
	sbbb	$0, %al
	testb	%al, %al
	je	.L1727
	leaq	.LC16(%rip), %rdi
	movl	$7, %ecx
	movq	%rdx, %rsi
	repz cmpsb
	seta	%al
	sbbb	$0, %al
	testb	%al, %al
	je	.L1727
	leaq	.LC14(%rip), %rdi
	movl	$8, %ecx
	movq	%rdx, %rsi
	repz cmpsb
	seta	%al
	sbbb	$0, %al
	testb	%al, %al
	je	.L1727
	leaq	.LC121(%rip), %rdi
	movl	$10, %ecx
	movq	%rdx, %rsi
	repz cmpsb
	seta	%al
	sbbb	$0, %al
	testb	%al, %al
	je	.L1727
	leaq	.LC122(%rip), %rdi
	movl	$12, %ecx
	movq	%rdx, %rsi
	repz cmpsb
	seta	%al
	sbbb	$0, %al
	testb	%al, %al
	je	.L1727
	leaq	928(%r15), %rdi
	leaq	.LC123(%rip), %rsi
	movl	$1, 336(%r15)
	xorl	%eax, %eax
	xorl	%ebx, %ebx
	call	printfPQExpBuffer@PLT
	jmp	.L1622
	.p2align 4,,10
	.p2align 3
.L1962:
	leaq	.LC16(%rip), %rdi
	call	strdup@PLT
	testq	%rax, %rax
	movq	%rax, 168(%r15)
	je	.L1637
	.p2align 4,,10
	.p2align 3
.L1727:
	movq	256(%r15), %rdi
	call	sslVerifyProtocolVersion
	testb	%al, %al
	movl	%eax, %ebx
	je	.L1963
	movq	264(%r15), %rdi
	call	sslVerifyProtocolVersion
	testb	%al, %al
	movl	%eax, %ebx
	je	.L1964
	movq	256(%r15), %rbp
	movq	264(%r15), %r12
	testq	%rbp, %rbp
	sete	%dl
	testq	%r12, %r12
	sete	%al
	orb	%al, %dl
	jne	.L1729
	cmpb	$0, 0(%rbp)
	je	.L1729
	cmpb	$0, (%r12)
	jne	.L1965
.L1729:
	movq	232(%r15), %rdx
	testq	%rdx, %rdx
	je	.L1732
	leaq	.LC118(%rip), %rdi
	movl	$8, %ecx
	movq	%rdx, %rsi
	repz cmpsb
	leaq	.LC14(%rip), %rdi
	movl	$8, %ecx
	movq	%rdx, %rsi
	seta	%r8b
	sbbb	$0, %r8b
	repz cmpsb
	seta	%al
	sbbb	$0, %al
	testb	%r8b, %r8b
	movsbl	%al, %eax
	je	.L1733
	leaq	.LC16(%rip), %rdi
	movl	$7, %ecx
	movq	%rdx, %rsi
	repz cmpsb
	seta	%cl
	sbbb	$0, %cl
	testb	%cl, %cl
	je	.L1733
	testl	%eax, %eax
	leaq	928(%r15), %rdi
	jne	.L1966
.L1734:
	leaq	.LC128(%rip), %rsi
	movl	$1, 336(%r15)
	xorl	%eax, %eax
	xorl	%ebx, %ebx
	call	printfPQExpBuffer@PLT
	jmp	.L1622
	.p2align 4,,10
	.p2align 3
.L1641:
	jne	.L1754
	testb	%cl, %cl
	je	.L1649
.L1754:
	movq	8(%r15), %rbp
.L1741:
	movl	$1, 336(%r15)
	movzbl	0(%rbp), %eax
	movl	$1, %edx
	testb	%al, %al
	je	.L1645
	.p2align 4,,10
	.p2align 3
.L1647:
	cmpb	$44, %al
	sete	%al
	addq	$1, %rbp
	movzbl	%al, %eax
	addl	%eax, %edx
	movzbl	0(%rbp), %eax
	testb	%al, %al
	jne	.L1647
.L1645:
	leaq	928(%r15), %rdi
	leaq	.LC111(%rip), %rsi
	movl	%ebx, %ecx
	xorl	%eax, %eax
	xorl	%ebx, %ebx
	call	printfPQExpBuffer@PLT
	jmp	.L1622
.L1661:
	cmpl	$1, %r13d
	je	.L1967
	cmpl	%r13d, %ecx
	jne	.L1743
	testb	%al, %al
	je	.L1658
	.p2align 4,,10
	.p2align 3
.L1743:
	movq	24(%r15), %rax
.L1744:
	movl	$1, 336(%r15)
	movzbl	(%rax), %esi
	movl	$1, %edx
	testb	%sil, %sil
	je	.L1669
	.p2align 4,,10
	.p2align 3
.L1671:
	cmpb	$44, %sil
	sete	%sil
	addq	$1, %rax
	movzbl	%sil, %esi
	addl	%esi, %edx
	movzbl	(%rax), %esi
	testb	%sil, %sil
	jne	.L1671
.L1669:
	leaq	928(%r15), %rdi
	leaq	.LC113(%rip), %rsi
	xorl	%eax, %eax
	xorl	%ebx, %ebx
	call	printfPQExpBuffer@PLT
	jmp	.L1622
	.p2align 4,,10
	.p2align 3
.L1963:
	movq	256(%r15), %rdx
	leaq	928(%r15), %rdi
	leaq	.LC124(%rip), %rsi
	movl	$1, 336(%r15)
	xorl	%eax, %eax
	call	printfPQExpBuffer@PLT
	jmp	.L1622
.L1964:
	movq	264(%r15), %rdx
	leaq	928(%r15), %rdi
	leaq	.LC125(%rip), %rsi
	movl	$1, 336(%r15)
	xorl	%eax, %eax
	call	printfPQExpBuffer@PLT
	jmp	.L1622
.L1961:
	leaq	.LC14(%rip), %rdi
	movl	$8, %ecx
	movq	%rdx, %rsi
	repz cmpsb
	seta	%al
	sbbb	$0, %al
	testb	%al, %al
	je	.L1723
	leaq	928(%r15), %rdi
	leaq	.LC119(%rip), %rsi
	movl	$1, 336(%r15)
	xorl	%eax, %eax
	xorl	%ebx, %ebx
	call	printfPQExpBuffer@PLT
	jmp	.L1622
.L1960:
	cmpb	$0, (%rax)
	je	.L1679
	jmp	.L1681
	.p2align 4,,10
	.p2align 3
.L1749:
	leaq	.LC110(%rip), %r13
	jmp	.L1691
.L1957:
	movq	32(%rsp), %rdx
	leaq	.LC115(%rip), %rsi
.L1947:
	movq	stderr@GOTPCREL(%rip), %rax
	movq	(%rax), %rdi
	xorl	%eax, %eax
	call	pg_fprintf@PLT
.L1946:
	xorl	%r12d, %r12d
.L1689:
	movq	8(%rsp), %rax
	addl	$1, 4(%rsp)
	movq	%r12, 32(%rax)
	movl	4(%rsp), %eax
	cmpl	%eax, 392(%r15)
	jg	.L1718
	jmp	.L1679
	.p2align 4,,10
	.p2align 3
.L1716:
	movq	%rbx, %rdi
	xorl	%r12d, %r12d
	call	fclose@PLT
	movl	$320, %edx
	movl	$320, %esi
	movq	%rbp, %rdi
	call	__explicit_bzero_chk@PLT
	jmp	.L1689
.L1956:
	movl	$1, 336(%r15)
	xorl	%ebx, %ebx
	jmp	.L1622
.L1751:
	leaq	.LC26(%rip), %rax
	movq	%rax, 40(%rsp)
	jmp	.L1692
.L1755:
	movq	%r12, %rsi
.L1711:
	movb	$0, (%rsi)
	jmp	.L1689
.L1954:
	cmpl	$1, %r13d
	jne	.L1743
	jmp	.L1667
	.p2align 4,,10
	.p2align 3
.L1968:
	addl	$1, %r13d
	cmpl	%r13d, 392(%r15)
	jle	.L1658
.L1667:
	movslq	%r13d, %rax
	movq	24(%r14), %rdi
	leaq	(%rax,%rax,4), %rbx
	salq	$3, %rbx
	call	strdup@PLT
	movq	%rax, 24(%r14,%rbx)
	movq	400(%r15), %r14
	cmpq	$0, 24(%r14,%rbx)
	jne	.L1968
	jmp	.L1637
.L1719:
	leaq	.LC16(%rip), %rdi
	call	strdup@PLT
	testq	%rax, %rax
	movq	%rax, 128(%r15)
	jne	.L1723
	jmp	.L1637
.L1733:
	testl	%eax, %eax
	je	.L1735
.L1738:
	movq	56(%r15), %rax
	testq	%rax, %rax
	je	.L1737
	leaq	.LC84(%rip), %rdi
	movl	$5, %ecx
	movq	%rax, %rsi
	repz cmpsb
	seta	%dl
	sbbb	$0, %dl
	testb	%dl, %dl
	je	.L1969
.L1737:
	movq	272(%r15), %rdx
	testq	%rdx, %rdx
	je	.L1740
	leaq	.LC129(%rip), %rdi
	movl	$4, %ecx
	movq	%rdx, %rsi
	repz cmpsb
	seta	%al
	sbbb	$0, %al
	testb	%al, %al
	je	.L1740
	leaq	.LC72(%rip), %rdi
	movl	$11, %ecx
	movq	%rdx, %rsi
	repz cmpsb
	seta	%al
	sbbb	$0, %al
	testb	%al, %al
	jne	.L1970
.L1740:
	movb	$1, 366(%r15)
	jmp	.L1622
.L1958:
	movq	32(%rsp), %rdx
	leaq	.LC116(%rip), %rsi
	jmp	.L1947
.L1732:
	leaq	.LC118(%rip), %rdi
	call	strdup@PLT
	testq	%rax, %rax
	movq	%rax, 232(%r15)
	jne	.L1738
	jmp	.L1637
.L1969:
	movq	%rax, %rdi
	call	free@PLT
	movl	$1, %esi
	xorl	%edi, %edi
	call	pg_get_encoding_from_locale@PLT
	movl	%eax, %edi
	call	pg_encoding_to_char@PLT
	movq	%rax, %rdi
	call	strdup@PLT
	testq	%rax, %rax
	movq	%rax, 56(%r15)
	jne	.L1737
	jmp	.L1637
.L1965:
	leaq	.LC22(%rip), %rsi
	movq	%rbp, %rdi
	call	pg_strcasecmp@PLT
	testl	%eax, %eax
	je	.L1729
	leaq	.LC22(%rip), %rsi
	movq	%r12, %rdi
	call	pg_strcasecmp@PLT
	testl	%eax, %eax
	jne	.L1730
.L1731:
	leaq	928(%r15), %rdi
	leaq	.LC126(%rip), %rsi
	movl	$1, 336(%r15)
	xorl	%eax, %eax
	xorl	%ebx, %ebx
	call	printfPQExpBuffer@PLT
	jmp	.L1622
.L1735:
	leaq	928(%r15), %rdi
	jmp	.L1734
.L1970:
	leaq	928(%r15), %rdi
	leaq	.LC130(%rip), %rsi
	movl	$1, 336(%r15)
	xorl	%eax, %eax
	xorl	%ebx, %ebx
	call	printfPQExpBuffer@PLT
	jmp	.L1622
.L1966:
	leaq	.LC127(%rip), %rsi
	movl	$1, 336(%r15)
	xorl	%eax, %eax
	xorl	%ebx, %ebx
	call	printfPQExpBuffer@PLT
	jmp	.L1622
.L1708:
	movl	$320, %edx
	movl	$320, %esi
	movq	%rbp, %rdi
	call	__explicit_bzero_chk@PLT
	jmp	.L1689
.L1730:
	movq	%r12, %rsi
	movq	%rbp, %rdi
	call	pg_strcasecmp@PLT
	testl	%eax, %eax
	jg	.L1731
	jmp	.L1729
.L1949:
	call	__stack_chk_fail@PLT
	.cfi_endproc
.LFE921:
	.size	connectOptions2, .-connectOptions2
	.p2align 4,,15
	.globl	PQconnectStartParams
	.type	PQconnectStartParams, @function
PQconnectStartParams:
.LFB915:
	.cfi_startproc
	pushq	%r15
	.cfi_def_cfa_offset 16
	.cfi_offset 15, -16
	pushq	%r14
	.cfi_def_cfa_offset 24
	.cfi_offset 14, -24
	movq	%rdi, %r15
	pushq	%r13
	.cfi_def_cfa_offset 32
	.cfi_offset 13, -32
	pushq	%r12
	.cfi_def_cfa_offset 40
	.cfi_offset 12, -40
	pushq	%rbp
	.cfi_def_cfa_offset 48
	.cfi_offset 6, -48
	pushq	%rbx
	.cfi_def_cfa_offset 56
	.cfi_offset 3, -56
	movl	%edx, %ebx
	subq	$56, %rsp
	.cfi_def_cfa_offset 112
	movq	%rdi, (%rsp)
	movq	%rsi, 8(%rsp)
	call	makeEmptyPGconn
	testq	%rax, %rax
	movq	%rax, 40(%rsp)
	je	.L1971
	addq	$928, %rax
	testl	%ebx, %ebx
	movq	%rax, 32(%rsp)
	je	.L1996
	movq	(%r15), %rsi
	testq	%rsi, %rsi
	je	.L1996
	leaq	.LC102(%rip), %r8
	xorl	%edx, %edx
	xorl	%ecx, %ecx
	.p2align 4,,10
	.p2align 3
.L1976:
	movq	8(%rsp), %rax
	movq	%r8, %rdi
	movq	(%rax,%rcx), %rbx
	movl	$7, %ecx
	repz cmpsb
	seta	%al
	sbbb	$0, %al
	testq	%rbx, %rbx
	movsbl	%al, %eax
	je	.L1974
	testl	%eax, %eax
	je	.L2053
.L1974:
	movq	(%rsp), %rsi
	addl	$1, %edx
	movslq	%edx, %rax
	leaq	0(,%rax,8), %rcx
	movq	(%rsi,%rax,8), %rsi
	testq	%rsi, %rsi
	jne	.L1976
.L1996:
	movq	$0, 24(%rsp)
.L1973:
	movq	32(%rsp), %rdi
	call	conninfo_init
	testq	%rax, %rax
	movq	%rax, %r12
	je	.L2052
	movq	(%rsp), %rax
	xorl	%r13d, %r13d
	xorl	%edx, %edx
	movq	(%rax), %r14
	testq	%r14, %r14
	jne	.L1978
	jmp	.L1979
	.p2align 4,,10
	.p2align 3
.L1980:
	movq	(%rsp), %rsi
	addl	$1, %r13d
	movslq	%r13d, %rax
	leaq	0(,%rax,8), %rdx
	movq	(%rsi,%rax,8), %r14
	testq	%r14, %r14
	je	.L1979
.L1978:
	movq	8(%rsp), %rax
	movq	(%rax,%rdx), %rbx
	testq	%rbx, %rbx
	je	.L1980
	cmpb	$0, (%rbx)
	je	.L1980
	movq	(%r12), %rdi
	testq	%rdi, %rdi
	je	.L1981
	movq	%r12, %rbp
	jmp	.L1983
	.p2align 4,,10
	.p2align 3
.L2054:
	addq	$56, %rbp
	movq	0(%rbp), %rdi
	testq	%rdi, %rdi
	je	.L1981
.L1983:
	movq	%r14, %rsi
	call	strcmp@PLT
	testl	%eax, %eax
	jne	.L2054
	leaq	.LC102(%rip), %rdi
	movl	$7, %ecx
	movq	%r14, %rsi
	repz cmpsb
	seta	%al
	sbbb	$0, %al
	testb	%al, %al
	jne	.L1993
	movq	24(%rsp), %r14
	testq	%r14, %r14
	jne	.L2050
.L1993:
	movq	24(%rbp), %rdi
	testq	%rdi, %rdi
	je	.L1989
	call	free@PLT
.L1989:
	movq	%rbx, %rdi
	call	strdup@PLT
	testq	%rax, %rax
	movq	%rax, 24(%rbp)
	jne	.L1980
.L2051:
	movq	32(%rsp), %rdi
	leaq	.LC3(%rip), %rsi
	call	printfPQExpBuffer@PLT
	movq	%r12, %rdi
	call	PQconninfoFree@PLT
.L2052:
	movq	24(%rsp), %rdi
	call	PQconninfoFree@PLT
	.p2align 4,,10
	.p2align 3
.L1975:
	movq	40(%rsp), %rax
	movl	$1, 336(%rax)
.L1971:
	movq	40(%rsp), %rax
	addq	$56, %rsp
	.cfi_remember_state
	.cfi_def_cfa_offset 56
	popq	%rbx
	.cfi_def_cfa_offset 48
	popq	%rbp
	.cfi_def_cfa_offset 40
	popq	%r12
	.cfi_def_cfa_offset 32
	popq	%r13
	.cfi_def_cfa_offset 24
	popq	%r14
	.cfi_def_cfa_offset 16
	popq	%r15
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L2056:
	.cfi_restore_state
	movq	24(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L1986
	call	free@PLT
	movq	24(%r14), %rax
	movq	%rax, 16(%rsp)
.L1986:
	movq	16(%rsp), %rdi
	call	strdup@PLT
	testq	%rax, %rax
	movq	%rax, 24(%rbx)
	je	.L2051
.L1984:
	addq	$56, %r14
.L2050:
	movq	(%r14), %rbp
	testq	%rbp, %rbp
	je	.L2055
	movq	24(%r14), %rax
	testq	%rax, %rax
	movq	%rax, 16(%rsp)
	je	.L1984
	movq	(%r12), %rdi
	testq	%rdi, %rdi
	je	.L1984
	movq	%r12, %rbx
	xorl	%r15d, %r15d
	jmp	.L1987
	.p2align 4,,10
	.p2align 3
.L1985:
	addl	$1, %r15d
	movslq	%r15d, %rcx
	leaq	0(,%rcx,8), %rax
	subq	%rcx, %rax
	leaq	(%r12,%rax,8), %rbx
	movq	(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L1984
.L1987:
	movq	%rbp, %rsi
	call	strcmp@PLT
	testl	%eax, %eax
	jne	.L1985
	jmp	.L2056
	.p2align 4,,10
	.p2align 3
.L1981:
	movq	32(%rsp), %rdi
	leaq	.LC18(%rip), %rsi
	movq	%r14, %rdx
	xorl	%eax, %eax
	call	printfPQExpBuffer@PLT
	movq	%r12, %rdi
	call	PQconninfoFree@PLT
	movq	24(%rsp), %rdi
	call	PQconninfoFree@PLT
	jmp	.L1975
	.p2align 4,,10
	.p2align 3
.L2055:
	movq	24(%rsp), %rdi
	call	PQconninfoFree@PLT
	movq	$0, 24(%rsp)
	jmp	.L1980
	.p2align 4,,10
	.p2align 3
.L1979:
	movq	24(%rsp), %rdi
	call	PQconninfoFree@PLT
	movq	32(%rsp), %rsi
	movq	%r12, %rdi
	call	conninfo_add_defaults
	testb	%al, %al
	je	.L2057
	movq	40(%rsp), %rdi
	movq	%r12, %rsi
	call	fillPGconn
	testb	%al, %al
	je	.L2058
	movq	%r12, %rdi
	call	PQconninfoFree@PLT
	movq	40(%rsp), %rdi
	call	connectOptions2
	testb	%al, %al
	je	.L1971
	movq	40(%rsp), %rbx
	movq	%rbx, %rdi
	call	connectDBStart
	testl	%eax, %eax
	jne	.L1971
	movl	$1, 336(%rbx)
	jmp	.L1971
.L2058:
	movq	%r12, %rdi
	call	PQconninfoFree@PLT
	jmp	.L1971
.L2053:
	movq	%rbx, %rdi
	call	recognized_connection_string
	testb	%al, %al
	movq	$0, 24(%rsp)
	je	.L1973
	movq	32(%rsp), %rsi
	xorl	%edx, %edx
	movq	%rbx, %rdi
	call	parse_connection_string
	testq	%rax, %rax
	movq	%rax, 24(%rsp)
	jne	.L1973
	jmp	.L1975
.L2057:
	movq	%r12, %rdi
	call	PQconninfoFree@PLT
	jmp	.L1975
	.cfi_endproc
.LFE915:
	.size	PQconnectStartParams, .-PQconnectStartParams
	.p2align 4,,15
	.globl	PQconnectdbParams
	.type	PQconnectdbParams, @function
PQconnectdbParams:
.LFB911:
	.cfi_startproc
	pushq	%rbx
	.cfi_def_cfa_offset 16
	.cfi_offset 3, -16
	call	PQconnectStartParams@PLT
	testq	%rax, %rax
	movq	%rax, %rbx
	je	.L2059
	cmpl	$1, 336(%rax)
	je	.L2059
	movq	%rax, %rdi
	call	connectDBComplete
.L2059:
	movq	%rbx, %rax
	popq	%rbx
	.cfi_def_cfa_offset 8
	ret
	.cfi_endproc
.LFE911:
	.size	PQconnectdbParams, .-PQconnectdbParams
	.p2align 4,,15
	.globl	PQpingParams
	.type	PQpingParams, @function
PQpingParams:
.LFB912:
	.cfi_startproc
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset 6, -16
	pushq	%rbx
	.cfi_def_cfa_offset 24
	.cfi_offset 3, -24
	subq	$8, %rsp
	.cfi_def_cfa_offset 32
	call	PQconnectStartParams@PLT
	movq	%rax, %rbx
	movq	%rax, %rdi
	call	internal_ping
	movq	%rbx, %rdi
	movl	%eax, %ebp
	call	PQfinish@PLT
	addq	$8, %rsp
	.cfi_def_cfa_offset 24
	movl	%ebp, %eax
	popq	%rbx
	.cfi_def_cfa_offset 16
	popq	%rbp
	.cfi_def_cfa_offset 8
	ret
	.cfi_endproc
.LFE912:
	.size	PQpingParams, .-PQpingParams
	.p2align 4,,15
	.globl	PQconnectStart
	.type	PQconnectStart, @function
PQconnectStart:
.LFB916:
	.cfi_startproc
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset 6, -16
	pushq	%rbx
	.cfi_def_cfa_offset 24
	.cfi_offset 3, -24
	movq	%rdi, %rbp
	subq	$8, %rsp
	.cfi_def_cfa_offset 32
	call	makeEmptyPGconn
	testq	%rax, %rax
	movq	%rax, %rbx
	je	.L2067
	movq	%rbp, %rsi
	movq	%rax, %rdi
	call	connectOptions1
	testb	%al, %al
	jne	.L2079
.L2067:
	addq	$8, %rsp
	.cfi_remember_state
	.cfi_def_cfa_offset 24
	movq	%rbx, %rax
	popq	%rbx
	.cfi_def_cfa_offset 16
	popq	%rbp
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L2079:
	.cfi_restore_state
	movq	%rbx, %rdi
	call	connectOptions2
	testb	%al, %al
	je	.L2067
	movq	%rbx, %rdi
	call	connectDBStart
	testl	%eax, %eax
	jne	.L2067
	movl	$1, 336(%rbx)
	jmp	.L2067
	.cfi_endproc
.LFE916:
	.size	PQconnectStart, .-PQconnectStart
	.p2align 4,,15
	.globl	PQconnectdb
	.type	PQconnectdb, @function
PQconnectdb:
.LFB913:
	.cfi_startproc
	pushq	%rbx
	.cfi_def_cfa_offset 16
	.cfi_offset 3, -16
	call	PQconnectStart@PLT
	testq	%rax, %rax
	movq	%rax, %rbx
	je	.L2080
	cmpl	$1, 336(%rax)
	je	.L2080
	movq	%rax, %rdi
	call	connectDBComplete
.L2080:
	movq	%rbx, %rax
	popq	%rbx
	.cfi_def_cfa_offset 8
	ret
	.cfi_endproc
.LFE913:
	.size	PQconnectdb, .-PQconnectdb
	.p2align 4,,15
	.globl	PQping
	.type	PQping, @function
PQping:
.LFB914:
	.cfi_startproc
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset 6, -16
	pushq	%rbx
	.cfi_def_cfa_offset 24
	.cfi_offset 3, -24
	subq	$8, %rsp
	.cfi_def_cfa_offset 32
	call	PQconnectStart@PLT
	movq	%rax, %rbx
	movq	%rax, %rdi
	call	internal_ping
	movq	%rbx, %rdi
	movl	%eax, %ebp
	call	PQfinish@PLT
	addq	$8, %rsp
	.cfi_def_cfa_offset 24
	movl	%ebp, %eax
	popq	%rbx
	.cfi_def_cfa_offset 16
	popq	%rbp
	.cfi_def_cfa_offset 8
	ret
	.cfi_endproc
.LFE914:
	.size	PQping, .-PQping
	.p2align 4,,15
	.globl	PQsetdbLogin
	.type	PQsetdbLogin, @function
PQsetdbLogin:
.LFB923:
	.cfi_startproc
	pushq	%r15
	.cfi_def_cfa_offset 16
	.cfi_offset 15, -16
	pushq	%r14
	.cfi_def_cfa_offset 24
	.cfi_offset 14, -24
	movq	%rcx, %r15
	pushq	%r13
	.cfi_def_cfa_offset 32
	.cfi_offset 13, -32
	pushq	%r12
	.cfi_def_cfa_offset 40
	.cfi_offset 12, -40
	movq	%rsi, %r13
	pushq	%rbp
	.cfi_def_cfa_offset 48
	.cfi_offset 6, -48
	pushq	%rbx
	.cfi_def_cfa_offset 56
	.cfi_offset 3, -56
	movq	%rdi, %r12
	movq	%rdx, %r14
	movq	%r8, %rbp
	subq	$24, %rsp
	.cfi_def_cfa_offset 80
	movq	%r9, 8(%rsp)
	call	makeEmptyPGconn
	testq	%rax, %rax
	movq	%rax, %rbx
	je	.L2088
	testq	%rbp, %rbp
	je	.L2090
	movq	%rbp, %rdi
	call	recognized_connection_string
	testb	%al, %al
	jne	.L2169
	leaq	.LC82(%rip), %rsi
	movq	%rbx, %rdi
	call	connectOptions1
	testb	%al, %al
	je	.L2088
	cmpb	$0, 0(%rbp)
	je	.L2094
	movq	88(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L2095
	call	free@PLT
.L2095:
	movq	%rbp, %rdi
	call	strdup@PLT
	testq	%rax, %rax
	movq	%rax, 88(%rbx)
	je	.L2096
.L2094:
	testq	%r12, %r12
	je	.L2093
	cmpb	$0, (%r12)
	je	.L2093
	movq	8(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L2099
	call	free@PLT
.L2099:
	movq	%r12, %rdi
	call	strdup@PLT
	testq	%rax, %rax
	movq	%rax, 8(%rbx)
	jne	.L2093
	.p2align 4,,10
	.p2align 3
.L2096:
	leaq	928(%rbx), %rdi
	leaq	.LC3(%rip), %rsi
	movl	$1, 336(%rbx)
	xorl	%eax, %eax
	call	printfPQExpBuffer@PLT
.L2088:
	addq	$24, %rsp
	.cfi_remember_state
	.cfi_def_cfa_offset 56
	movq	%rbx, %rax
	popq	%rbx
	.cfi_def_cfa_offset 48
	popq	%rbp
	.cfi_def_cfa_offset 40
	popq	%r12
	.cfi_def_cfa_offset 32
	popq	%r13
	.cfi_def_cfa_offset 24
	popq	%r14
	.cfi_def_cfa_offset 16
	popq	%r15
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L2093:
	.cfi_restore_state
	testq	%r13, %r13
	je	.L2098
	cmpb	$0, 0(%r13)
	je	.L2098
	movq	24(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L2102
	call	free@PLT
.L2102:
	movq	%r13, %rdi
	call	strdup@PLT
	testq	%rax, %rax
	movq	%rax, 24(%rbx)
	je	.L2096
.L2098:
	testq	%r14, %r14
	je	.L2101
	cmpb	$0, (%r14)
	je	.L2101
	movq	64(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L2105
	call	free@PLT
.L2105:
	movq	%r14, %rdi
	call	strdup@PLT
	testq	%rax, %rax
	movq	%rax, 64(%rbx)
	je	.L2096
.L2101:
	testq	%r15, %r15
	je	.L2104
	cmpb	$0, (%r15)
	je	.L2104
	movq	32(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L2108
	call	free@PLT
.L2108:
	movq	%r15, %rdi
	call	strdup@PLT
	testq	%rax, %rax
	movq	%rax, 32(%rbx)
	je	.L2096
.L2104:
	cmpq	$0, 8(%rsp)
	je	.L2107
	movq	8(%rsp), %rax
	cmpb	$0, (%rax)
	je	.L2107
	movq	104(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L2111
	call	free@PLT
.L2111:
	movq	8(%rsp), %rdi
	call	strdup@PLT
	testq	%rax, %rax
	movq	%rax, 104(%rbx)
	je	.L2096
.L2107:
	cmpq	$0, 80(%rsp)
	je	.L2110
	movq	80(%rsp), %rax
	cmpb	$0, (%rax)
	je	.L2110
	movq	112(%rbx), %rdi
	testq	%rdi, %rdi
	je	.L2113
	call	free@PLT
.L2113:
	movq	80(%rsp), %rdi
	call	strdup@PLT
	testq	%rax, %rax
	movq	%rax, 112(%rbx)
	je	.L2096
.L2110:
	movq	%rbx, %rdi
	call	connectOptions2
	testb	%al, %al
	je	.L2088
	movq	%rbx, %rdi
	call	connectDBStart
	testl	%eax, %eax
	je	.L2088
	movq	%rbx, %rdi
	call	connectDBComplete
	jmp	.L2088
	.p2align 4,,10
	.p2align 3
.L2169:
	movq	%rbp, %rsi
	movq	%rbx, %rdi
	call	connectOptions1
	testb	%al, %al
	jne	.L2094
	jmp	.L2088
	.p2align 4,,10
	.p2align 3
.L2090:
	leaq	.LC82(%rip), %rsi
	movq	%rax, %rdi
	call	connectOptions1
	testb	%al, %al
	jne	.L2094
	jmp	.L2088
	.cfi_endproc
.LFE923:
	.size	PQsetdbLogin, .-PQsetdbLogin
	.p2align 4,,15
	.globl	PQregisterThreadLock
	.type	PQregisterThreadLock, @function
PQregisterThreadLock:
.LFB1009:
	.cfi_startproc
	movq	pg_g_threadlock@GOTPCREL(%rip), %rcx
	leaq	default_threadlock(%rip), %rdx
	testq	%rdi, %rdi
	cmove	%rdx, %rdi
	movq	(%rcx), %rax
	movq	%rdi, (%rcx)
	ret
	.cfi_endproc
.LFE1009:
	.size	PQregisterThreadLock, .-PQregisterThreadLock
	.section	.rodata
	.align 16
	.type	query.25836, @object
	.size	query.25836, 28
query.25836:
	.string	"set client_encoding to '%s'"
	.local	singlethread_lock.25935
	.comm	singlethread_lock.25935,40,32
	.globl	pg_g_threadlock
	.section	.data.rel.local,"aw",@progbits
	.align 8
	.type	pg_g_threadlock, @object
	.size	pg_g_threadlock, 8
pg_g_threadlock:
	.quad	default_threadlock
	.section	.rodata
	.align 8
	.type	short_uri_designator, @object
	.size	short_uri_designator, 12
short_uri_designator:
	.string	"postgres://"
	.align 8
	.type	uri_designator, @object
	.size	uri_designator, 14
uri_designator:
	.string	"postgresql://"
	.section	.rodata.str1.1
.LC131:
	.string	"PGDATESTYLE"
.LC132:
	.string	"datestyle"
.LC133:
	.string	"PGTZ"
.LC134:
	.string	"timezone"
.LC135:
	.string	"PGGEQO"
.LC136:
	.string	"geqo"
	.section	.data.rel.ro.local,"aw",@progbits
	.align 32
	.type	EnvironmentOptions, @object
	.size	EnvironmentOptions, 64
EnvironmentOptions:
	.quad	.LC131
	.quad	.LC132
	.quad	.LC133
	.quad	.LC134
	.quad	.LC135
	.quad	.LC136
	.quad	0
	.quad	0
	.section	.rodata.str1.1
.LC137:
	.string	"authtype"
.LC138:
	.string	"PGAUTHTYPE"
.LC139:
	.string	"Database-Authtype"
.LC140:
	.string	"D"
.LC141:
	.string	"Database-Service"
.LC142:
	.string	"usedbauth"
.LC143:
	.string	"USEDBAUTH"
.LC144:
	.string	"no"
.LC145:
	.string	"Use DBAuth"
.LC146:
	.string	"PGUSER"
.LC147:
	.string	"Database-User"
.LC148:
	.string	"PGPASSWORD"
.LC149:
	.string	"Database-Password"
.LC150:
	.string	"*"
.LC151:
	.string	"passfile"
.LC152:
	.string	"PGPASSFILE"
.LC153:
	.string	"Database-Password-File"
.LC154:
	.string	"channel_binding"
.LC155:
	.string	"PGCHANNELBINDING"
.LC156:
	.string	"Channel-Binding"
.LC157:
	.string	"PGCONNECT_TIMEOUT"
.LC158:
	.string	"Connect-timeout"
.LC159:
	.string	"PGDATABASE"
.LC160:
	.string	"Database-Name"
.LC161:
	.string	"PGHOST"
.LC162:
	.string	"Database-Host"
.LC163:
	.string	"hostaddr"
.LC164:
	.string	"PGHOSTADDR"
.LC165:
	.string	"Database-Host-IP-Address"
.LC166:
	.string	"PGPORT"
.LC167:
	.string	"Database-Port"
.LC168:
	.string	"PGCLIENTENCODING"
.LC169:
	.string	"Client-Encoding"
.LC170:
	.string	"tty"
.LC171:
	.string	"PGTTY"
.LC172:
	.string	"Backend-Debug-TTY"
.LC173:
	.string	"options"
.LC174:
	.string	"PGOPTIONS"
.LC175:
	.string	"Backend-Options"
.LC176:
	.string	"application_name"
.LC177:
	.string	"PGAPPNAME"
.LC178:
	.string	"Application-Name"
.LC179:
	.string	"fallback_application_name"
.LC180:
	.string	"Fallback-Application-Name"
.LC181:
	.string	"keepalives"
.LC182:
	.string	"TCP-Keepalives"
.LC183:
	.string	"TCP-Keepalives-Idle"
.LC184:
	.string	"TCP-Keepalives-Interval"
.LC185:
	.string	"TCP-Keepalives-Count"
.LC186:
	.string	"TCP-User-Timeout"
.LC187:
	.string	"PGSSLMODE"
.LC188:
	.string	"SSL-Mode"
.LC189:
	.string	"sslcompression"
.LC190:
	.string	"PGSSLCOMPRESSION"
.LC191:
	.string	"0"
.LC192:
	.string	"SSL-Compression"
.LC193:
	.string	"sslcert"
.LC194:
	.string	"PGSSLCERT"
.LC195:
	.string	"SSL-Client-Cert"
.LC196:
	.string	"sslkey"
.LC197:
	.string	"PGSSLKEY"
.LC198:
	.string	"SSL-Client-Key"
.LC199:
	.string	"sslpassword"
.LC200:
	.string	"SSL-Client-Key-Password"
.LC201:
	.string	"sslrootcert"
.LC202:
	.string	"PGSSLROOTCERT"
.LC203:
	.string	"SSL-Root-Certificate"
.LC204:
	.string	"sslcrl"
.LC205:
	.string	"PGSSLCRL"
.LC206:
	.string	"SSL-Revocation-List"
.LC207:
	.string	"requirepeer"
.LC208:
	.string	"PGREQUIREPEER"
.LC209:
	.string	"Require-Peer"
.LC210:
	.string	"ssl_min_protocol_version"
.LC211:
	.string	"PGSSLMINPROTOCOLVERSION"
.LC212:
	.string	"SSL-Minimum-Protocol-Version"
.LC213:
	.string	"ssl_max_protocol_version"
.LC214:
	.string	"PGSSLMAXPROTOCOLVERSION"
.LC215:
	.string	"SSL-Maximum-Protocol-Version"
.LC216:
	.string	"gssencmode"
.LC217:
	.string	"PGGSSENCMODE"
.LC218:
	.string	"GSSENC-Mode"
.LC219:
	.string	"krbsrvname"
.LC220:
	.string	"PGKRBSRVNAME"
.LC221:
	.string	"postgres"
.LC222:
	.string	"Kerberos-service-name"
.LC223:
	.string	"gsslib"
.LC224:
	.string	"PGGSSLIB"
.LC225:
	.string	"GSS-library"
.LC226:
	.string	"replication"
.LC227:
	.string	"Replication"
.LC228:
	.string	"target_session_attrs"
.LC229:
	.string	"PGTARGETSESSIONATTRS"
.LC230:
	.string	"Target-Session-Attrs"
	.section	.data.rel.ro.local
	.align 32
	.type	PQconninfoOptions, @object
	.size	PQconninfoOptions, 2432
PQconninfoOptions:
	.quad	.LC137
	.quad	.LC138
	.quad	.LC82
	.quad	0
	.quad	.LC139
	.quad	.LC140
	.long	20
	.zero	4
	.quad	-1
	.quad	.LC2
	.quad	.LC87
	.quad	0
	.quad	0
	.quad	.LC141
	.quad	.LC82
	.long	20
	.zero	4
	.quad	-1
	.quad	.LC142
	.quad	.LC143
	.quad	.LC144
	.quad	0
	.quad	.LC145
	.quad	.LC82
	.long	20
	.zero	4
	.quad	0
	.quad	.LC95
	.quad	.LC146
	.quad	0
	.quad	0
	.quad	.LC147
	.quad	.LC82
	.long	20
	.zero	4
	.quad	104
	.quad	.LC97
	.quad	.LC148
	.quad	0
	.quad	0
	.quad	.LC149
	.quad	.LC150
	.long	20
	.zero	4
	.quad	112
	.quad	.LC151
	.quad	.LC152
	.quad	0
	.quad	0
	.quad	.LC153
	.quad	.LC82
	.long	64
	.zero	4
	.quad	120
	.quad	.LC154
	.quad	.LC155
	.quad	0
	.quad	0
	.quad	.LC156
	.quad	.LC82
	.long	8
	.zero	4
	.quad	128
	.quad	.LC78
	.quad	.LC157
	.quad	0
	.quad	0
	.quad	.LC158
	.quad	.LC82
	.long	10
	.zero	4
	.quad	40
	.quad	.LC102
	.quad	.LC159
	.quad	0
	.quad	0
	.quad	.LC160
	.quad	.LC82
	.long	20
	.zero	4
	.quad	88
	.quad	.LC101
	.quad	.LC161
	.quad	0
	.quad	0
	.quad	.LC162
	.quad	.LC82
	.long	40
	.zero	4
	.quad	8
	.quad	.LC163
	.quad	.LC164
	.quad	0
	.quad	0
	.quad	.LC165
	.quad	.LC82
	.long	45
	.zero	4
	.quad	16
	.quad	.LC33
	.quad	.LC166
	.quad	.LC26
	.quad	0
	.quad	.LC167
	.quad	.LC82
	.long	6
	.zero	4
	.quad	24
	.quad	.LC85
	.quad	.LC168
	.quad	0
	.quad	0
	.quad	.LC169
	.quad	.LC82
	.long	10
	.zero	4
	.quad	56
	.quad	.LC170
	.quad	.LC171
	.quad	.LC82
	.quad	0
	.quad	.LC172
	.quad	.LC140
	.long	40
	.zero	4
	.quad	32
	.quad	.LC173
	.quad	.LC174
	.quad	.LC82
	.quad	0
	.quad	.LC175
	.quad	.LC82
	.long	40
	.zero	4
	.quad	64
	.quad	.LC176
	.quad	.LC177
	.quad	0
	.quad	0
	.quad	.LC178
	.quad	.LC82
	.long	64
	.zero	4
	.quad	72
	.quad	.LC179
	.quad	0
	.quad	0
	.quad	0
	.quad	.LC180
	.quad	.LC82
	.long	64
	.zero	4
	.quad	80
	.quad	.LC181
	.quad	0
	.quad	0
	.quad	0
	.quad	.LC182
	.quad	.LC82
	.long	1
	.zero	4
	.quad	136
	.quad	.LC48
	.quad	0
	.quad	0
	.quad	0
	.quad	.LC183
	.quad	.LC82
	.long	10
	.zero	4
	.quad	144
	.quad	.LC50
	.quad	0
	.quad	0
	.quad	0
	.quad	.LC184
	.quad	.LC82
	.long	10
	.zero	4
	.quad	152
	.quad	.LC52
	.quad	0
	.quad	0
	.quad	0
	.quad	.LC185
	.quad	.LC82
	.long	10
	.zero	4
	.quad	160
	.quad	.LC54
	.quad	0
	.quad	0
	.quad	0
	.quad	.LC186
	.quad	.LC82
	.long	10
	.zero	4
	.quad	48
	.quad	.LC15
	.quad	.LC187
	.quad	.LC16
	.quad	0
	.quad	.LC188
	.quad	.LC82
	.long	12
	.zero	4
	.quad	168
	.quad	.LC189
	.quad	.LC190
	.quad	.LC191
	.quad	0
	.quad	.LC192
	.quad	.LC82
	.long	1
	.zero	4
	.quad	176
	.quad	.LC193
	.quad	.LC194
	.quad	0
	.quad	0
	.quad	.LC195
	.quad	.LC82
	.long	64
	.zero	4
	.quad	192
	.quad	.LC196
	.quad	.LC197
	.quad	0
	.quad	0
	.quad	.LC198
	.quad	.LC82
	.long	64
	.zero	4
	.quad	184
	.quad	.LC199
	.quad	0
	.quad	0
	.quad	0
	.quad	.LC200
	.quad	.LC150
	.long	20
	.zero	4
	.quad	200
	.quad	.LC201
	.quad	.LC202
	.quad	0
	.quad	0
	.quad	.LC203
	.quad	.LC82
	.long	64
	.zero	4
	.quad	208
	.quad	.LC204
	.quad	.LC205
	.quad	0
	.quad	0
	.quad	.LC206
	.quad	.LC82
	.long	64
	.zero	4
	.quad	216
	.quad	.LC207
	.quad	.LC208
	.quad	0
	.quad	0
	.quad	.LC209
	.quad	.LC82
	.long	10
	.zero	4
	.quad	224
	.quad	.LC210
	.quad	.LC211
	.quad	0
	.quad	0
	.quad	.LC212
	.quad	.LC82
	.long	8
	.zero	4
	.quad	256
	.quad	.LC213
	.quad	.LC214
	.quad	0
	.quad	0
	.quad	.LC215
	.quad	.LC82
	.long	8
	.zero	4
	.quad	264
	.quad	.LC216
	.quad	.LC217
	.quad	.LC118
	.quad	0
	.quad	.LC218
	.quad	.LC82
	.long	8
	.zero	4
	.quad	232
	.quad	.LC219
	.quad	.LC220
	.quad	.LC221
	.quad	0
	.quad	.LC222
	.quad	.LC82
	.long	20
	.zero	4
	.quad	240
	.quad	.LC223
	.quad	.LC224
	.quad	0
	.quad	0
	.quad	.LC225
	.quad	.LC82
	.long	7
	.zero	4
	.quad	248
	.quad	.LC226
	.quad	0
	.quad	0
	.quad	0
	.quad	.LC227
	.quad	.LC140
	.long	5
	.zero	4
	.quad	96
	.quad	.LC228
	.quad	.LC229
	.quad	.LC129
	.quad	0
	.quad	.LC230
	.quad	.LC82
	.long	11
	.zero	4
	.quad	272
	.quad	0
	.quad	0
	.quad	0
	.quad	0
	.quad	0
	.quad	0
	.long	0
	.zero	12
	.ident	"GCC: (Ubuntu 7.5.0-3ubuntu1~18.04) 7.5.0"
	.section	.note.GNU-stack,"",@progbits
