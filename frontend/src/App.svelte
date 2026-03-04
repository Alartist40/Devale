<script lang="ts">
    import { onMount } from 'svelte';
    import { EventsOn } from '../wailsjs/runtime/runtime.js';
    import { RunCommand, GetSystemInfo, RunRepairPhase, ScheduleResume, ClearResume, SaveState, LoadState } from '../wailsjs/go/main/App.js';

    let currentTab = 'Home';
    let logs: string[] = [];
    let commandInput = '';
    let sysInfo = { cpu: 'Loading...', memory: 'Loading...', gpu: 'Loading...', disk: 'Loading...' };
    let repairPhase = 0;
    let isRepairing = false;

    onMount(async () => {
        EventsOn('terminal:output', (line: string) => {
            logs = [...logs, line];
            // Auto-scroll
            setTimeout(() => {
                const el = document.getElementById('terminal-logs');
                if (el) el.scrollTop = el.scrollHeight;
            }, 10);
        });

        try {
            sysInfo = await GetSystemInfo();
            const lastPhase = await LoadState();
            if (lastPhase > 0) {
                currentTab = 'Home';
                // If we rebooted at phase 2, we should resume at phase 3
                repairPhase = lastPhase === 2 ? 3 : lastPhase;
                logs = [...logs, `>>> RESUMING REPAIR AT PHASE ${repairPhase}...`];
                continueRepair();
            }
        } catch (e) {
            console.error(e);
        }
    });

    async function handleCommand(e: KeyboardEvent) {
        if (e.key === 'Enter' && commandInput.trim()) {
            const cmd = commandInput;
            commandInput = '';
            logs = [...logs, `USER@DEVALE > ${cmd}`];
            await RunCommand(cmd);
        }
    }

    async function continueRepair() {
        isRepairing = true;
        for (let i = repairPhase || 1; i <= 6; i++) {
            repairPhase = i;
            await SaveState(i);

            if (i === 2) {
                logs = [...logs, "!!! PHASE 2 REQUIRES RESTART !!!"];
                await RunRepairPhase(2);
                await ScheduleResume();
                logs = [...logs, "System will restart in 10 seconds. Auto-resume scheduled."];
                await RunCommand("shutdown /r /t 10 /c \"DevAle System Repair - Restarting for CHKDSK\"");
                return;
            }

            await RunRepairPhase(i);
        }

        await ClearResume();
        await SaveState(0);
        repairPhase = 0;
        isRepairing = false;
        logs = [...logs, "--- ALL REPAIR PHASES COMPLETED SUCCESSFULLY ---"];
    }

    async function startFullRepair() {
        if (isRepairing) return;
        logs = [...logs, "--- INITIATING FULL SYSTEM REPAIR (6 PHASES) ---"];
        repairPhase = 1;
        continueRepair();
    }
</script>

<div class="sidebar">
    <div style="padding: 0 20px 20px; font-size: 24px; font-weight: bold; color: var(--accent-magenta); text-shadow: 0 0 10px var(--accent-magenta);">DEVALE v2</div>
    <div role="button" tabindex="0" class="nav-item {currentTab === 'Home' ? 'active' : ''}" on:click={() => currentTab = 'Home'} on:keydown={(e) => e.key === 'Enter' && (currentTab = 'Home')}>HOME</div>
    <div role="button" tabindex="0" class="nav-item {currentTab === 'Diagnose' ? 'active' : ''}" on:click={() => currentTab = 'Diagnose'} on:keydown={(e) => e.key === 'Enter' && (currentTab = 'Diagnose')}>DIAGNOSE</div>
    <div role="button" tabindex="0" class="nav-item {currentTab === 'Tools' ? 'active' : ''}" on:click={() => currentTab = 'Tools'} on:keydown={(e) => e.key === 'Enter' && (currentTab = 'Tools')}>TOOLS</div>
    <div role="button" tabindex="0" class="nav-item {currentTab === 'AppStore' ? 'active' : ''}" on:click={() => currentTab = 'AppStore'} on:keydown={(e) => e.key === 'Enter' && (currentTab = 'AppStore')}>APP STORE</div>
</div>

<main class="main-content">
    <div class="content-area">
        {#if currentTab === 'Home'}
            <div style="text-align: center;">
                <h1 style="margin-bottom: 5px;">SYSTEM STATUS: {isRepairing ? 'REPAIRING...' : 'READY'}</h1>
                <div style="color: var(--accent-cyan); font-family: monospace; margin-bottom: 20px;">[ PHASE {repairPhase}/6 ]</div>

                <button class="btn-panic" on:click={startFullRepair} disabled={isRepairing}>
                    {isRepairing ? 'RUNNING...' : 'PANIC\nBUTTON'}
                </button>

                <div style="margin-top: 20px; display: flex; justify-content: center; gap: 10px;">
                    {#each [1,2,3,4,5,6] as p}
                        <div style="width: 30px; height: 4px; background: {repairPhase >= p ? 'var(--accent-cyan)' : '#333'}; box-shadow: {repairPhase >= p ? '0 0 10px var(--accent-cyan)' : 'none'}"></div>
                    {/each}
                </div>
                <p style="margin-top: 30px; color: var(--text-dim); max-width: 400px; margin-left: auto; margin-right: auto;">
                    Automated deep-system repair. Includes DISM, SFC, CHKDSK, WMI, and Store App fixes. May require a restart.
                </p>
            </div>
        {:else if currentTab === 'Diagnose'}
            <h1>Hardware Telemetry</h1>
            <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 20px;">
                <div class="card">
                    <h3 style="color: var(--accent-magenta)">CPU PROCEESOR</h3>
                    <p style="font-family: monospace;">{sysInfo.cpu}</p>
                </div>
                <div class="card">
                    <h3 style="color: var(--accent-magenta)">PHYSICAL MEMORY</h3>
                    <p style="font-family: monospace;">{sysInfo.memory}</p>
                </div>
                <div class="card">
                    <h3 style="color: var(--accent-magenta)">GRAPHICS ADAPTER</h3>
                    <p style="font-family: monospace;">{sysInfo.gpu}</p>
                </div>
                <div class="card">
                    <h3 style="color: var(--accent-magenta)">BLOCK STORAGE</h3>
                    <p style="font-family: monospace;">{sysInfo.disk}</p>
                </div>
            </div>
        {:else if currentTab === 'Tools'}
            <h1>Manual Overrides</h1>
            <div style="display: flex; gap: 15px; flex-wrap: wrap;">
                <button class="tool-btn" on:click={() => RunCommand("ipconfig /flushdns")}>FLUSH_DNS</button>
                <button class="tool-btn" on:click={() => RunCommand("taskkill /f /im explorer.exe && start explorer.exe")}>RESTART_EXPLORER</button>
                <button class="tool-btn" on:click={() => RunCommand("cleanmgr /sagerun:1")}>DISK_CLEANUP</button>
                <button class="tool-btn" on:click={() => RunCommand("SystemPropertiesProtection.exe")}>RESTORE_POINT_UI</button>
            </div>
        {:else if currentTab === 'AppStore'}
            <h1>NEURAL APP STORE</h1>
            <div style="background: #111; padding: 20px; border: 1px solid var(--accent-cyan);">
                <p style="color: var(--accent-cyan)">> winget integration active</p>
                <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 10px;">
                    <div class="app-item">
                        <span>Google Chrome</span>
                        <button on:click={() => RunCommand("winget install Google.Chrome")}>INSTALL</button>
                    </div>
                    <div class="app-item">
                        <span>VS Code</span>
                        <button on:click={() => RunCommand("winget install Microsoft.VisualStudioCode")}>INSTALL</button>
                    </div>
                    <div class="app-item">
                        <span>Discord</span>
                        <button on:click={() => RunCommand("winget install Discord.Discord")}>INSTALL</button>
                    </div>
                    <div class="app-item">
                        <span>7-Zip</span>
                        <button on:click={() => RunCommand("winget install 7zip.7zip")}>INSTALL</button>
                    </div>
                </div>
            </div>
        {/if}
    </div>

    <div class="terminal-area">
        <div id="terminal-logs" class="terminal-logs">
            {#each logs as log}
                <div>{log}</div>
            {/each}
        </div>
        <div class="terminal-input-row">
            <span class="terminal-prompt">USER@DEVALE ></span>
            <input class="terminal-input" bind:value={commandInput} on:keydown={handleCommand} placeholder="TYPE COMMAND..." />
        </div>
    </div>
</main>

<style>
    .card {
        background: #16161a;
        padding: 20px;
        border: 1px solid var(--border);
        box-shadow: inset 0 0 15px rgba(0,0,0,0.5);
    }
    .tool-btn {
        padding: 12px 24px;
        background: transparent;
        border: 1px solid var(--accent-cyan);
        color: var(--accent-cyan);
        cursor: pointer;
        font-family: monospace;
        letter-spacing: 1px;
    }
    .tool-btn:hover {
        background: var(--accent-cyan);
        color: black;
        box-shadow: 0 0 15px var(--accent-cyan);
    }
    .app-item {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 10px;
        border-bottom: 1px solid #222;
    }
    .app-item button {
        background: transparent;
        border: 1px solid var(--accent-magenta);
        color: var(--accent-magenta);
        padding: 5px 10px;
        cursor: pointer;
    }
    .app-item button:hover {
        background: var(--accent-magenta);
        color: white;
    }
</style>
