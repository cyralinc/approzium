	.file	"fe-auth.c"
	.text
	.section	.rodata.str1.1,"aMS",@progbits,1
.LC0:
	.string	"SCRAM-SHA-256-PLUS"
.LC1:
	.string	"SCRAM-SHA-256"
	.section	.rodata.str1.8,"aMS",@progbits,1
	.align 8
.LC2:
	.string	"channel binding required, but server authenticated client without channel binding\n"
	.align 8
.LC3:
	.string	"channel binding required but not supported by server's authentication request\n"
	.align 8
.LC4:
	.string	"Kerberos 4 authentication not supported\n"
	.align 8
.LC5:
	.string	"Kerberos 5 authentication not supported\n"
	.align 8
.LC6:
	.string	"GSSAPI authentication not supported\n"
	.align 8
.LC7:
	.string	"SSPI authentication not supported\n"
	.align 8
.LC8:
	.string	"Crypt authentication not supported\n"
	.align 8
.LC9:
	.string	"fe_sendauth: no password supplied\n"
	.section	.rodata.str1.1
.LC10:
	.string	"out of memory\n"
	.section	.rodata.str1.8
	.align 8
.LC11:
	.string	"fe_sendauth: error sending password authentication\n"
	.align 8
.LC12:
	.string	"channel binding required, but SSL not in use\n"
	.align 8
.LC13:
	.string	"duplicate SASL authentication request\n"
	.align 8
.LC14:
	.string	"fe_sendauth: invalid authentication request from server: invalid list of authentication mechanisms\n"
	.align 8
.LC15:
	.string	"server offered SCRAM-SHA-256-PLUS authentication over a non-SSL connection\n"
	.align 8
.LC16:
	.string	"none of the server's SASL authentication mechanisms are supported\n"
	.align 8
.LC17:
	.string	"channel binding is required, but server did not offer an authentication method that supports channel binding\n"
	.align 8
.LC18:
	.string	"fe_sendauth: invalid authentication request from server: AUTH_REQ_SASL_CONT without AUTH_REQ_SASL\n"
	.align 8
.LC19:
	.string	"out of memory allocating SASL buffer (%d)\n"
	.align 8
.LC20:
	.string	"AuthenticationSASLFinal received from server, but SASL authentication was not completed\n"
	.align 8
.LC21:
	.string	"fe_sendauth: error in SASL authentication\n"
	.align 8
.LC22:
	.string	"SCM_CRED authentication method not supported\n"
	.align 8
.LC23:
	.string	"authentication method %u not supported\n"
	.text
	.p2align 4,,15
	.globl	pg_fe_sendauth
	.type	pg_fe_sendauth, @function
pg_fe_sendauth:
.LFB900:
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
	movl	%esi, %r12d
	pushq	%rbp
	.cfi_def_cfa_offset 48
	.cfi_offset 6, -48
	pushq	%rbx
	.cfi_def_cfa_offset 56
	.cfi_offset 3, -56
	movl	%edi, %ebp
	movq	%rdx, %rbx
	subq	$72, %rsp
	.cfi_def_cfa_offset 128
	movq	%fs:40, %rax
	movq	%rax, 56(%rsp)
	xorl	%eax, %eax
	movq	120(%rdx), %rax
	cmpb	$114, (%rax)
	je	.L112
.L2:
	cmpl	$12, %ebp
	ja	.L7
	leaq	.L8(%rip), %rcx
	movl	%ebp, %edx
	movslq	(%rcx,%rdx,4), %rax
	addq	%rcx, %rax
	jmp	*%rax
	.section	.rodata
	.align 4
	.align 4
.L8:
	.long	.L68-.L8
	.long	.L9-.L8
	.long	.L10-.L8
	.long	.L11-.L8
	.long	.L12-.L8
	.long	.L11-.L8
	.long	.L13-.L8
	.long	.L14-.L8
	.long	.L14-.L8
	.long	.L15-.L8
	.long	.L16-.L8
	.long	.L17-.L8
	.long	.L17-.L8
	.text
	.p2align 4,,10
	.p2align 3
.L112:
	testl	%edi, %edi
	je	.L4
	leal	-10(%rdi), %eax
	cmpl	$2, %eax
	jbe	.L2
	leaq	920(%rdx), %rdi
	leaq	.LC3(%rip), %rsi
	xorl	%eax, %eax
	movl	$-1, %r13d
	call	printfPQExpBuffer@PLT
	.p2align 4,,10
	.p2align 3
.L1:
	movq	56(%rsp), %rsi
	xorq	%fs:40, %rsi
	movl	%r13d, %eax
	jne	.L113
	addq	$72, %rsp
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
.L17:
	.cfi_restore_state
	cmpq	$0, 880(%rbx)
	je	.L114
	leal	1(%r12), %edi
	movslq	%edi, %rdi
	call	malloc@PLT
	testq	%rax, %rax
	movq	%rax, %r14
	je	.L115
	movslq	%r12d, %r15
	movq	%rbx, %rdx
	movq	%rax, %rdi
	movq	%r15, %rsi
	call	pqGetnchar@PLT
	testl	%eax, %eax
	movl	%eax, %r13d
	jne	.L116
	leaq	16(%rsp), %rcx
	subq	$8, %rsp
	.cfi_def_cfa_offset 136
	movq	880(%rbx), %rdi
	leaq	12(%rsp), %rax
	movb	$0, (%r14,%r15)
	movl	%r12d, %edx
	movq	%r14, %rsi
	pushq	%rax
	.cfi_def_cfa_offset 144
	leaq	19(%rsp), %r9
	leaq	24(%rsp), %r8
	call	pg_fe_scram_exchange@PLT
	movq	%r14, %rdi
	call	free@PLT
	popq	%rax
	.cfi_def_cfa_offset 136
	cmpl	$12, %ebp
	popq	%rdx
	.cfi_def_cfa_offset 128
	movslq	8(%rsp), %rcx
	je	.L117
	testl	%ecx, %ecx
	jne	.L62
.L66:
	cmpb	$0, 3(%rsp)
	je	.L1
.L63:
	cmpb	$0, 4(%rsp)
	jne	.L1
.L57:
	cmpq	$0, 928(%rbx)
	movl	$-1, %r13d
	jne	.L1
	leaq	920(%rbx), %rdi
	leaq	.LC21(%rip), %rsi
	xorl	%eax, %eax
	call	printfPQExpBuffer@PLT
	jmp	.L1
	.p2align 4,,10
	.p2align 3
.L9:
	leaq	920(%rbx), %rdi
	leaq	.LC4(%rip), %rsi
	xorl	%eax, %eax
	movl	$-1, %r13d
	call	printfPQExpBuffer@PLT
	jmp	.L1
	.p2align 4,,10
	.p2align 3
.L10:
	leaq	920(%rbx), %rdi
	leaq	.LC5(%rip), %rsi
	xorl	%eax, %eax
	movl	$-1, %r13d
	call	printfPQExpBuffer@PLT
	jmp	.L1
	.p2align 4,,10
	.p2align 3
.L11:
	movslq	388(%rbx), %rax
	movq	392(%rbx), %r8
	movb	$1, 697(%rbx)
	leaq	(%rax,%rax,4), %rax
	leaq	(%r8,%rax,8), %rax
	movq	32(%rax), %r12
	testq	%r12, %r12
	je	.L118
.L18:
	cmpb	$0, (%r12)
	je	.L19
	cmpl	$5, %ebp
	je	.L119
	cmpl	$3, %ebp
	jne	.L109
	xorl	%ebp, %ebp
.L25:
	movq	%r12, %rdi
	call	strlen@PLT
	cmpw	$2, 690(%rbx)
	leaq	1(%rax), %rcx
	movq	%r12, %rdx
	jbe	.L28
	movl	$112, %esi
	movq	%rbx, %rdi
	call	pqPacketSend@PLT
	movl	%eax, %r12d
.L29:
	testq	%rbp, %rbp
	je	.L30
	movq	%rbp, %rdi
	call	free@PLT
.L30:
	testl	%r12d, %r12d
	jne	.L109
.L68:
	xorl	%r13d, %r13d
	jmp	.L1
	.p2align 4,,10
	.p2align 3
.L12:
	leaq	920(%rbx), %rdi
	leaq	.LC8(%rip), %rsi
	xorl	%eax, %eax
	movl	$-1, %r13d
	call	printfPQExpBuffer@PLT
	jmp	.L1
	.p2align 4,,10
	.p2align 3
.L13:
	leaq	920(%rbx), %rdi
	leaq	.LC22(%rip), %rsi
	xorl	%eax, %eax
	movl	$-1, %r13d
	call	printfPQExpBuffer@PLT
	jmp	.L1
	.p2align 4,,10
	.p2align 3
.L14:
	leaq	920(%rbx), %rdi
	leaq	.LC6(%rip), %rsi
	xorl	%eax, %eax
	movl	$-1, %r13d
	call	printfPQExpBuffer@PLT
	jmp	.L1
	.p2align 4,,10
	.p2align 3
.L15:
	leaq	920(%rbx), %rdi
	leaq	.LC7(%rip), %rsi
	xorl	%eax, %eax
	movl	$-1, %r13d
	call	printfPQExpBuffer@PLT
	jmp	.L1
	.p2align 4,,10
	.p2align 3
.L16:
	leaq	16(%rsp), %rbp
	movq	$0, 8(%rsp)
	movq	%rbp, %rdi
	call	initPQExpBuffer@PLT
	movq	120(%rbx), %rax
	cmpb	$114, (%rax)
	je	.L120
.L32:
	cmpq	$0, 880(%rbx)
	je	.L71
	leaq	920(%rbx), %rdi
	leaq	.LC13(%rip), %rsi
	xorl	%eax, %eax
	call	printfPQExpBuffer@PLT
.L33:
	movq	%rbp, %rdi
	movl	$-1, %r13d
	call	termPQExpBuffer@PLT
	movq	8(%rsp), %rdi
	testq	%rdi, %rdi
	je	.L1
.L110:
	call	free@PLT
	jmp	.L1
	.p2align 4,,10
	.p2align 3
.L4:
	movq	880(%rdx), %rdi
	call	pg_fe_scram_channel_bound@PLT
	testb	%al, %al
	jne	.L68
	leaq	920(%rbx), %rdi
	leaq	.LC2(%rip), %rsi
	xorl	%eax, %eax
	movl	$-1, %r13d
	call	printfPQExpBuffer@PLT
	jmp	.L1
	.p2align 4,,10
	.p2align 3
.L7:
	leaq	920(%rbx), %rdi
	leaq	.LC23(%rip), %rsi
	movl	%ebp, %edx
	xorl	%eax, %eax
	movl	$-1, %r13d
	call	printfPQExpBuffer@PLT
	jmp	.L1
	.p2align 4,,10
	.p2align 3
.L117:
	cmpb	$0, 3(%rsp)
	je	.L121
	testl	%ecx, %ecx
	je	.L63
.L62:
	movq	16(%rsp), %rdx
	movl	$112, %esi
	movq	%rbx, %rdi
	call	pqPacketSend@PLT
	movq	16(%rsp), %rdi
	movl	%eax, %ebp
	call	free@PLT
	testl	%ebp, %ebp
	je	.L66
	jmp	.L57
	.p2align 4,,10
	.p2align 3
.L71:
	leaq	.LC0(%rip), %r12
	leaq	.LC1(%rip), %r13
	xorl	%r14d, %r14d
.L34:
	movq	%rbx, %rsi
	movq	%rbp, %rdi
	call	pqGets@PLT
	testl	%eax, %eax
	jne	.L122
.L35:
	cmpq	$0, 32(%rsp)
	je	.L36
	movq	16(%rsp), %rax
	cmpb	$0, (%rax)
	je	.L37
	movl	$19, %ecx
	movq	%rax, %rsi
	movq	%r12, %rdi
	repz cmpsb
	seta	%dl
	sbbb	$0, %dl
	testb	%dl, %dl
	jne	.L38
	cmpb	$0, 888(%rbx)
	je	.L39
	movq	120(%rbx), %rax
	movq	%rbx, %rsi
	movq	%rbp, %rdi
	cmpb	$100, (%rax)
	cmovne	%r12, %r14
	call	pqGets@PLT
	testl	%eax, %eax
	je	.L35
.L122:
	leaq	920(%rbx), %rdi
	leaq	.LC14(%rip), %rsi
	xorl	%eax, %eax
	call	printfPQExpBuffer@PLT
	jmp	.L33
	.p2align 4,,10
	.p2align 3
.L109:
	addq	$920, %rbx
.L23:
	leaq	.LC11(%rip), %rsi
	movq	%rbx, %rdi
	xorl	%eax, %eax
	orl	$-1, %r13d
	call	printfPQExpBuffer@PLT
	jmp	.L1
	.p2align 4,,10
	.p2align 3
.L38:
	movq	%rax, %rsi
	movl	$14, %ecx
	movq	%r13, %rdi
	repz cmpsb
	seta	%al
	sbbb	$0, %al
	testq	%r14, %r14
	movsbl	%al, %eax
	jne	.L34
	testl	%eax, %eax
	cmove	%r13, %r14
	jmp	.L34
	.p2align 4,,10
	.p2align 3
.L120:
	cmpb	$0, 888(%rbx)
	jne	.L32
	leaq	920(%rbx), %rdi
	leaq	.LC12(%rip), %rsi
	xorl	%eax, %eax
	call	printfPQExpBuffer@PLT
	jmp	.L33
	.p2align 4,,10
	.p2align 3
.L28:
	xorl	%esi, %esi
	movq	%rbx, %rdi
	call	pqPacketSend@PLT
	movl	%eax, %r12d
	jmp	.L29
	.p2align 4,,10
	.p2align 3
.L118:
	movq	104(%rbx), %r12
	testq	%r12, %r12
	jne	.L18
.L19:
	leaq	920(%rbx), %rdi
	leaq	.LC9(%rip), %rsi
	xorl	%eax, %eax
	movl	$-1, %r13d
	call	printfPQExpBuffer@PLT
	jmp	.L1
	.p2align 4,,10
	.p2align 3
.L119:
	leaq	52(%rsp), %r14
	movq	%rbx, %rdx
	movl	$4, %esi
	movq	%r14, %rdi
	call	pqGetnchar@PLT
	testl	%eax, %eax
	jne	.L109
	movl	$72, %edi
	call	malloc@PLT
	testq	%rax, %rax
	movq	%rax, %rbp
	je	.L123
	movq	96(%rbx), %r13
	movq	%r13, %rdi
	call	strlen@PLT
	leaq	36(%rbp), %rcx
	movq	%rax, %rdx
	movq	%r13, %rsi
	movq	%r12, %rdi
	call	pg_md5_encrypt@PLT
	testb	%al, %al
	je	.L108
	leaq	39(%rbp), %rdi
	movq	%rbp, %rcx
	movl	$4, %edx
	movq	%r14, %rsi
	call	pg_md5_encrypt@PLT
	testb	%al, %al
	je	.L108
	movq	%rbp, %r12
	jmp	.L25
	.p2align 4,,10
	.p2align 3
.L36:
	movq	%rbp, %rdi
	call	termPQExpBuffer@PLT
	movq	8(%rsp), %rdi
	testq	%rdi, %rdi
	je	.L54
	call	free@PLT
.L54:
	leaq	920(%rbx), %rdi
	leaq	.LC10(%rip), %rsi
	xorl	%eax, %eax
	movl	$-1, %r13d
	call	printfPQExpBuffer@PLT
	jmp	.L1
	.p2align 4,,10
	.p2align 3
.L37:
	testq	%r14, %r14
	je	.L124
	movq	120(%rbx), %rax
	cmpb	$114, (%rax)
	je	.L125
.L43:
	movslq	388(%rbx), %rax
	movq	392(%rbx), %rdi
	movb	$1, 697(%rbx)
	leaq	(%rax,%rax,4), %rax
	leaq	(%rdi,%rax,8), %rax
	movq	32(%rax), %rsi
	testq	%rsi, %rsi
	je	.L126
.L44:
	cmpb	$0, (%rsi)
	jne	.L46
.L45:
	leaq	920(%rbx), %rdi
	leaq	.LC9(%rip), %rsi
	xorl	%eax, %eax
	call	printfPQExpBuffer@PLT
	jmp	.L33
.L39:
	leaq	920(%rbx), %rdi
	leaq	.LC15(%rip), %rsi
	xorl	%eax, %eax
	call	printfPQExpBuffer@PLT
	jmp	.L33
.L125:
	leaq	.LC0(%rip), %rdi
	movl	$19, %ecx
	movq	%r14, %rsi
	repz cmpsb
	seta	%al
	sbbb	$0, %al
	testb	%al, %al
	je	.L43
	leaq	920(%rbx), %rdi
	leaq	.LC17(%rip), %rsi
	xorl	%eax, %eax
	call	printfPQExpBuffer@PLT
	jmp	.L33
.L126:
	movq	104(%rbx), %rsi
	testq	%rsi, %rsi
	je	.L45
	jmp	.L44
.L46:
	movq	%r14, %rdx
	movq	%rbx, %rdi
	call	pg_fe_scram_init@PLT
	testq	%rax, %rax
	movq	%rax, 880(%rbx)
	je	.L36
	leaq	8(%rsp), %rcx
	subq	$8, %rsp
	.cfi_def_cfa_offset 136
	xorl	%esi, %esi
	leaq	11(%rsp), %rdx
	movq	%rax, %rdi
	pushq	%rdx
	.cfi_def_cfa_offset 144
	movl	$-1, %edx
	leaq	18(%rsp), %r9
	leaq	20(%rsp), %r8
	call	pg_fe_scram_exchange@PLT
	popq	%rcx
	.cfi_def_cfa_offset 136
	popq	%rsi
	.cfi_def_cfa_offset 128
	cmpb	$0, 2(%rsp)
	je	.L47
	cmpb	$0, 3(%rsp)
	je	.L33
.L47:
	movq	%rbx, %rdx
	movl	$1, %esi
	movl	$112, %edi
	call	pqPutMsgStart@PLT
	testl	%eax, %eax
	jne	.L33
	movq	%rbx, %rsi
	movq	%r14, %rdi
	call	pqPuts@PLT
	testl	%eax, %eax
	jne	.L33
	cmpq	$0, 8(%rsp)
	je	.L51
	movl	4(%rsp), %edi
	movq	%rbx, %rdx
	movl	$4, %esi
	call	pqPutInt@PLT
	testl	%eax, %eax
	jne	.L33
	movslq	4(%rsp), %rsi
	movq	8(%rsp), %rdi
	movq	%rbx, %rdx
	call	pqPutnchar@PLT
	testl	%eax, %eax
	jne	.L33
.L51:
	movq	%rbx, %rdi
	call	pqPutMsgEnd@PLT
	testl	%eax, %eax
	jne	.L33
	movq	%rbx, %rdi
	call	pqFlush@PLT
	testl	%eax, %eax
	movl	%eax, %r13d
	jne	.L33
	movq	%rbp, %rdi
	call	termPQExpBuffer@PLT
	movq	8(%rsp), %rdi
	testq	%rdi, %rdi
	jne	.L110
	jmp	.L1
	.p2align 4,,10
	.p2align 3
.L116:
	movq	%r14, %rdi
	call	free@PLT
	jmp	.L57
.L124:
	leaq	920(%rbx), %rdi
	leaq	.LC16(%rip), %rsi
	xorl	%eax, %eax
	call	printfPQExpBuffer@PLT
	jmp	.L33
.L121:
	testl	%ecx, %ecx
	je	.L61
	movq	16(%rsp), %rdi
	call	free@PLT
.L61:
	leaq	920(%rbx), %rdi
	leaq	.LC20(%rip), %rsi
	xorl	%eax, %eax
	call	printfPQExpBuffer@PLT
	jmp	.L57
.L108:
	movq	%rbp, %rdi
	addq	$920, %rbx
	call	free@PLT
	jmp	.L23
.L115:
	leaq	920(%rbx), %rdi
	leaq	.LC19(%rip), %rsi
	movl	%r12d, %edx
	xorl	%eax, %eax
	call	printfPQExpBuffer@PLT
	jmp	.L57
.L114:
	leaq	920(%rbx), %rdi
	leaq	.LC18(%rip), %rsi
	xorl	%eax, %eax
	orl	$-1, %r13d
	call	printfPQExpBuffer@PLT
	jmp	.L1
.L113:
	call	__stack_chk_fail@PLT
.L123:
	addq	$920, %rbx
	leaq	.LC10(%rip), %rsi
	xorl	%eax, %eax
	movq	%rbx, %rdi
	call	printfPQExpBuffer@PLT
	jmp	.L23
	.cfi_endproc
.LFE900:
	.size	pg_fe_sendauth, .-pg_fe_sendauth
	.section	.rodata.str1.8
	.align 8
.LC24:
	.string	"could not look up local user ID %d: %s\n"
	.align 8
.LC25:
	.string	"local user with ID %d does not exist\n"
	.text
	.p2align 4,,15
	.globl	pg_fe_getauthname
	.type	pg_fe_getauthname, @function
pg_fe_getauthname:
.LFB901:
	.cfi_startproc
	pushq	%r13
	.cfi_def_cfa_offset 16
	.cfi_offset 13, -16
	pushq	%r12
	.cfi_def_cfa_offset 24
	.cfi_offset 12, -24
	movq	%rdi, %r12
	pushq	%rbp
	.cfi_def_cfa_offset 32
	.cfi_offset 6, -32
	pushq	%rbx
	.cfi_def_cfa_offset 40
	.cfi_offset 3, -40
	subq	$8280, %rsp
	.cfi_def_cfa_offset 8320
	movq	%fs:40, %rax
	movq	%rax, 8264(%rsp)
	xorl	%eax, %eax
	call	geteuid@PLT
	movq	pg_g_threadlock@GOTPCREL(%rip), %rbx
	leaq	64(%rsp), %r13
	movl	%eax, %ebp
	movl	$1, %edi
	movq	$0, 8(%rsp)
	call	*(%rbx)
	leaq	16(%rsp), %rsi
	leaq	8(%rsp), %r8
	movq	%r13, %rdx
	movl	$8192, %ecx
	movl	%ebp, %edi
	call	pqGetpwuid@PLT
	movq	8(%rsp), %rdx
	testq	%rdx, %rdx
	je	.L128
	movq	(%rdx), %rdi
	testq	%rdi, %rdi
	je	.L133
	call	strdup@PLT
	testq	%rax, %rax
	movq	%rax, %rbp
	jne	.L130
	testq	%r12, %r12
	je	.L130
	leaq	.LC10(%rip), %rsi
	movq	%r12, %rdi
	xorl	%eax, %eax
	call	printfPQExpBuffer@PLT
.L130:
	xorl	%edi, %edi
	call	*(%rbx)
	movq	8264(%rsp), %rcx
	xorq	%fs:40, %rcx
	movq	%rbp, %rax
	jne	.L141
	addq	$8280, %rsp
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
.L128:
	.cfi_restore_state
	testq	%r12, %r12
	je	.L133
	testl	%eax, %eax
	je	.L131
	movl	$8192, %edx
	movq	%r13, %rsi
	movl	%eax, %edi
	call	pg_strerror_r@PLT
	leaq	.LC24(%rip), %rsi
	movq	%rax, %rcx
	movl	%ebp, %edx
	movq	%r12, %rdi
	xorl	%eax, %eax
	xorl	%ebp, %ebp
	call	printfPQExpBuffer@PLT
	jmp	.L130
	.p2align 4,,10
	.p2align 3
.L131:
	leaq	.LC25(%rip), %rsi
	movl	%ebp, %edx
	movq	%r12, %rdi
	xorl	%eax, %eax
	xorl	%ebp, %ebp
	call	printfPQExpBuffer@PLT
	jmp	.L130
	.p2align 4,,10
	.p2align 3
.L133:
	xorl	%ebp, %ebp
	jmp	.L130
.L141:
	call	__stack_chk_fail@PLT
	.cfi_endproc
.LFE901:
	.size	pg_fe_getauthname, .-pg_fe_getauthname
	.p2align 4,,15
	.globl	PQencryptPassword
	.type	PQencryptPassword, @function
PQencryptPassword:
.LFB902:
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
	movl	$36, %edi
	movq	%rsi, %rbp
	call	malloc@PLT
	testq	%rax, %rax
	movq	%rax, %rbx
	je	.L142
	movq	%rbp, %rdi
	call	strlen@PLT
	movq	%rbx, %rcx
	movq	%rax, %rdx
	movq	%rbp, %rsi
	movq	%r12, %rdi
	call	pg_md5_encrypt@PLT
	testb	%al, %al
	je	.L148
.L142:
	movq	%rbx, %rax
	popq	%rbx
	.cfi_remember_state
	.cfi_def_cfa_offset 24
	popq	%rbp
	.cfi_def_cfa_offset 16
	popq	%r12
	.cfi_def_cfa_offset 8
	ret
	.p2align 4,,10
	.p2align 3
.L148:
	.cfi_restore_state
	movq	%rbx, %rdi
	xorl	%ebx, %ebx
	call	free@PLT
	jmp	.L142
	.cfi_endproc
.LFE902:
	.size	PQencryptPassword, .-PQencryptPassword
	.section	.rodata.str1.1
.LC26:
	.string	"show password_encryption"
	.section	.rodata.str1.8
	.align 8
.LC27:
	.string	"unexpected shape of result set returned for SHOW\n"
	.align 8
.LC28:
	.string	"password_encryption value too long\n"
	.section	.rodata.str1.1
.LC29:
	.string	"on"
.LC30:
	.string	"off"
.LC31:
	.string	"scram-sha-256"
.LC32:
	.string	"md5"
	.section	.rodata.str1.8
	.align 8
.LC33:
	.string	"unrecognized password encryption algorithm \"%s\"\n"
	.text
	.p2align 4,,15
	.globl	PQencryptPasswordConn
	.type	PQencryptPasswordConn, @function
PQencryptPasswordConn:
.LFB903:
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
	subq	$88, %rsp
	.cfi_def_cfa_offset 144
	movq	%fs:40, %rax
	movq	%rax, 72(%rsp)
	xorl	%eax, %eax
	testq	%rdi, %rdi
	je	.L180
	testq	%rcx, %rcx
	movq	%rdx, %r12
	movq	%rdi, %rbx
	movq	%rsi, %rbp
	movq	%rcx, %rdx
	je	.L181
.L152:
	leaq	.LC29(%rip), %rdi
	movl	$3, %ecx
	movq	%rdx, %rsi
	repz cmpsb
	seta	%al
	sbbb	$0, %al
	testb	%al, %al
	je	.L158
	leaq	.LC30(%rip), %rdi
	movl	$4, %ecx
	movq	%rdx, %rsi
	repz cmpsb
	seta	%al
	sbbb	$0, %al
	testb	%al, %al
	je	.L158
	leaq	.LC31(%rip), %rdi
	movl	$14, %ecx
	movq	%rdx, %rsi
	repz cmpsb
	seta	%al
	sbbb	$0, %al
	testb	%al, %al
	jne	.L159
	movq	%rbp, %rdi
	call	pg_fe_scram_build_secret@PLT
	testq	%rax, %rax
	movq	%rax, %r13
	je	.L162
.L149:
	movq	72(%rsp), %rbx
	xorq	%fs:40, %rbx
	movq	%r13, %rax
	jne	.L182
	addq	$88, %rsp
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
.L158:
	.cfi_restore_state
	movl	$36, %edi
	call	malloc@PLT
	testq	%rax, %rax
	movq	%rax, %r14
	movq	%rax, %r13
	je	.L162
	movq	%r12, %rdi
	call	strlen@PLT
	movq	%r14, %rcx
	movq	%rax, %rdx
	movq	%r12, %rsi
	movq	%rbp, %rdi
	call	pg_md5_encrypt@PLT
	testb	%al, %al
	jne	.L149
	movq	%r14, %rdi
	call	free@PLT
.L162:
	leaq	920(%rbx), %rdi
	leaq	.LC10(%rip), %rsi
	xorl	%eax, %eax
	call	printfPQExpBuffer@PLT
.L180:
	xorl	%r13d, %r13d
	jmp	.L149
	.p2align 4,,10
	.p2align 3
.L181:
	leaq	.LC26(%rip), %rsi
	call	PQexec@PLT
	testq	%rax, %rax
	movq	%rax, %r13
	je	.L180
	movq	%rax, %rdi
	call	PQresultStatus@PLT
	cmpl	$2, %eax
	movq	%r13, %rdi
	jne	.L183
	call	PQntuples@PLT
	cmpl	$1, %eax
	jne	.L156
	movq	%r13, %rdi
	call	PQnfields@PLT
	cmpl	$1, %eax
	jne	.L156
	xorl	%edx, %edx
	xorl	%esi, %esi
	movq	%r13, %rdi
	call	PQgetvalue@PLT
	movq	%rax, %rdi
	movq	%rax, %r14
	call	strlen@PLT
	cmpq	$50, %rax
	ja	.L184
	leaq	16(%rsp), %r15
	leaq	1(%rax), %rdx
	movl	$51, %ecx
	movq	%r14, %rsi
	movq	%r15, %rdi
	call	__memcpy_chk@PLT
	movq	%r13, %rdi
	call	PQclear@PLT
	movq	%r15, %rdx
	jmp	.L152
	.p2align 4,,10
	.p2align 3
.L156:
	movq	%r13, %rdi
	xorl	%r13d, %r13d
	call	PQclear@PLT
	leaq	920(%rbx), %rdi
	leaq	.LC27(%rip), %rsi
	xorl	%eax, %eax
	call	printfPQExpBuffer@PLT
	jmp	.L149
	.p2align 4,,10
	.p2align 3
.L183:
	call	PQclear@PLT
	xorl	%r13d, %r13d
	jmp	.L149
.L184:
	movq	%r13, %rdi
	xorl	%r13d, %r13d
	call	PQclear@PLT
	leaq	920(%rbx), %rdi
	leaq	.LC28(%rip), %rsi
	xorl	%eax, %eax
	call	printfPQExpBuffer@PLT
	jmp	.L149
.L182:
	call	__stack_chk_fail@PLT
.L159:
	leaq	.LC32(%rip), %rsi
	movq	%rdx, %rdi
	movq	%rdx, 8(%rsp)
	call	strcmp@PLT
	testl	%eax, %eax
	movq	8(%rsp), %rdx
	je	.L158
	leaq	920(%rbx), %rdi
	leaq	.LC33(%rip), %rsi
	xorl	%eax, %eax
	xorl	%r13d, %r13d
	call	printfPQExpBuffer@PLT
	jmp	.L149
	.cfi_endproc
.LFE903:
	.size	PQencryptPasswordConn, .-PQencryptPasswordConn
	.ident	"GCC: (Ubuntu 7.5.0-3ubuntu1~18.04) 7.5.0"
	.section	.note.GNU-stack,"",@progbits
