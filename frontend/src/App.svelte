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
    let theme = 'light';

    onMount(async () => {
        EventsOn('terminal:output', (line: string) => {
            logs = [...logs, line];
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
                repairPhase = lastPhase === 2 ? 3 : lastPhase;
                logs = [...logs, `>>> RESUMING REPAIR AT PHASE ${repairPhase}...`];
                void continueRepair();
            }
        } catch (e) {
            console.error(e);
        }
    });

    function toggleTheme() {
        theme = theme === 'light' ? 'dark' : 'light';
        document.body.setAttribute('data-theme', theme);
    }

    async function handleCommand(e: KeyboardEvent) {
        if (e.key === 'Enter' && commandInput.trim()) {
            const cmd = commandInput;
            commandInput = '';
            logs = [...logs, `USER@DEVALE > ${cmd}`];
            try {
                const res = await RunCommand(cmd);
                if (res.startsWith("Error:")) {
                    logs = [...logs, res];
                }
            } catch (err) {
                logs = [...logs, `ERROR: ${err instanceof Error ? err.message : String(err)}`];
            }
        }
    }

    async function continueRepair() {
        isRepairing = true;
        try {
            for (let i = repairPhase || 1; i <= 6; i++) {
                repairPhase = i;
                await SaveState(i);

                if (i === 2) {
                    logs = [...logs, "!!! PHASE 2 REQUIRES RESTART !!!"];
                    const phase2Result = await RunRepairPhase(2);
                    if (phase2Result.startsWith("Error:")) {
                        throw new Error(`Phase 2 failed: ${phase2Result}`);
                    }
                    const scheduleResult = await ScheduleResume();
                    if (scheduleResult.startsWith("Error:")) {
                        throw new Error(`Resume scheduling failed: ${scheduleResult}`);
                    }
                    logs = [...logs, "System will restart in 10 seconds. Auto-resume scheduled."];
                    await RunCommand("shutdown /r /t 10 /c \"DevAle System Repair - Restarting for CHKDSK\"");
                    return;
                }

                const result = await RunRepairPhase(i);
                if (result.startsWith("Error:")) {
                    throw new Error(`Phase ${i} failed: ${result}`);
                }
            }

            await ClearResume();
            await SaveState(0);
            repairPhase = 0;
            logs = [...logs, "--- ALL REPAIR PHASES COMPLETED SUCCESSFULLY ---"];
        } catch (e) {
            logs = [...logs, `ERROR: ${e instanceof Error ? e.message : String(e)}`];
        } finally {
            isRepairing = false;
        }
    }

    async function resetRepairState() {
        repairPhase = 0;
        isRepairing = false;
        await SaveState(0);
        await ClearResume();
        logs = [...logs, ">>> REPAIR STATE MANUALLY RESET"];
    }

    async function startFullRepair() {
        if (isRepairing) return;
        logs = [...logs, "--- INITIATING FULL SYSTEM REPAIR (6 PHASES) ---"];
        repairPhase = 1;
        void continueRepair();
    }
</script>

<div class="theme-switch-container" on:click={toggleTheme} role="button" tabindex="0" on:keydown={(e) => e.key === 'Enter' && toggleTheme()}>
    <span class="switch-label">{theme === 'light' ? 'LIGHT' : 'DARK'}</span>
    <div class="switch-base">
        <div class="switch-thumb {theme === 'dark' ? 'active' : ''}"></div>
    </div>
</div>

<div class="sidebar">
    <div style="padding: 0 20px 30px; font-size: 26px; font-weight: 900; color: var(--accent); text-align: center;">DEVALE</div>

    <div role="button" tabindex="0" class="nav-item {currentTab === 'Home' ? 'active' : ''}" on:click={() => currentTab = 'Home'} on:keydown={(e) => e.key === 'Enter' && (currentTab = 'Home')}>HOME</div>
    <div role="button" tabindex="0" class="nav-item {currentTab === 'Diagnose' ? 'active' : ''}" on:click={() => currentTab = 'Diagnose'} on:keydown={(e) => e.key === 'Enter' && (currentTab = 'Diagnose')}>DIAGNOSE</div>
    <div role="button" tabindex="0" class="nav-item {currentTab === 'Tools' ? 'active' : ''}" on:click={() => currentTab = 'Tools'} on:keydown={(e) => e.key === 'Enter' && (currentTab = 'Tools')}>TOOLS</div>
    <div role="button" tabindex="0" class="nav-item {currentTab === 'AppStore' ? 'active' : ''}" on:click={() => currentTab = 'AppStore'} on:keydown={(e) => e.key === 'Enter' && (currentTab = 'AppStore')}>APP STORE</div>
</div>

<main class="main-content">
    <div class="content-area">
        {#if currentTab === 'Home'}
            <div style="text-align: center;">
                <h1 style="margin-bottom: 5px; font-size: 2.5rem; font-weight: 900;">{isRepairing ? 'REPAIRING...' : 'SYSTEM READY'}</h1>
                <div style="color: var(--accent); font-family: monospace; font-weight: bold; margin-bottom: 20px;">[ PHASE {repairPhase}/6 ]</div>

                <div style="position: relative; width: 220px; margin: 40px auto;">
                    <button class="btn-panic" class:active={isRepairing} on:click={startFullRepair} disabled={isRepairing}>
                        PANIC<br>BUTTON
                    </button>

                    {#if !isRepairing && repairPhase > 0}
                        <button
                            on:click={resetRepairState}
                            style="position: absolute; bottom: -30px; left: 50%; transform: translateX(-50%); background: none; border: none; color: var(--text-dim); font-size: 10px; cursor: pointer; text-decoration: underline;"
                        >
                            RESET REPAIR STATE
                        </button>
                    {/if}
                </div>

                <div style="margin-top: 30px; display: flex; justify-content: center; gap: 15px;">
                    {#each [1,2,3,4,5,6] as p}
                        <div style="width: 40px; height: 6px; border-radius: 3px; background: var(--bg-main); box-shadow: {repairPhase >= p ? 'inset 2px 2px 4px var(--shadow-dark), inset -2px -2px 4px var(--shadow-light)' : '2px 2px 4px var(--shadow-dark), -2px -2px 4px var(--shadow-light)'};">
                            <div style="width: 100%; height: 100%; background: {repairPhase >= p ? 'var(--accent)' : 'transparent'}; border-radius: 3px; transition: background 0.3s;"></div>
                        </div>
                    {/each}
                </div>
            </div>
        {:else if currentTab === 'Diagnose'}
            <h1 style="font-weight: 900;">HARDWARE TELEMETRY</h1>
            <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 30px;">
                <div class="card">
                    <h3 style="color: var(--accent); margin-top: 0;">PROCESSOR</h3>
                    <p style="font-family: monospace; font-size: 1.1rem;">{sysInfo.cpu}</p>
                </div>
                <div class="card">
                    <h3 style="color: var(--accent); margin-top: 0;">MEMORY</h3>
                    <p style="font-family: monospace; font-size: 1.1rem;">{sysInfo.memory}</p>
                </div>
                <div class="card">
                    <h3 style="color: var(--accent); margin-top: 0;">GRAPHICS</h3>
                    <p style="font-family: monospace; font-size: 1.1rem;">{sysInfo.gpu}</p>
                </div>
                <div class="card">
                    <h3 style="color: var(--accent); margin-top: 0;">STORAGE</h3>
                    <p style="font-family: monospace; font-size: 1.1rem;">{sysInfo.disk}</p>
                </div>
            </div>
        {:else if currentTab === 'Tools'}
            <h1 style="font-weight: 900;">MANUAL OVERRIDES</h1>
            <div class="card" style="display: flex; gap: 15px; flex-wrap: wrap; justify-content: center;">
                <button class="tool-btn" on:click={() => RunCommand("ipconfig /flushdns")}>FLUSH_DNS</button>
                <button class="tool-btn" on:click={() => RunCommand("taskkill /f /im explorer.exe && start explorer.exe")}>RESTART_EXPLORER</button>
                <button class="tool-btn" on:click={() => RunCommand("cleanmgr /sagerun:1")}>DISK_CLEANUP</button>
                <button class="tool-btn" on:click={() => RunCommand("SystemPropertiesProtection.exe")}>RESTORE_POINT</button>
            </div>
        {:else if currentTab === 'AppStore'}
            <h1 style="font-weight: 900;">APP STORE</h1>
            <div class="card">
                <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 20px;">
                    {#each ["Google.Chrome", "Microsoft.VisualStudioCode", "Discord.Discord", "7zip.7zip"] as id}
                        <div class="app-item" style="padding: 15px; background: var(--bg-main); border-radius: 50px; box-shadow: 4px 4px 8px var(--shadow-dark), -4px -4px 8px var(--shadow-light); display: flex; justify-content: space-between; align-items: center;">
                            <span style="font-weight: bold; margin-left: 10px;">{id.split('.')[id.split('.').length-1]}</span>
                            <button class="tool-btn" style="padding: 8px 15px; font-size: 12px; margin: 0;" on:click={() => RunCommand(`winget install ${id}`)}>INSTALL</button>
                        </div>
                    {/each}
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
            <span class="terminal-prompt">SYSTEM ></span>
            <input class="terminal-input" bind:value={commandInput} on:keydown={handleCommand} placeholder="ENTER COMMAND..." />
        </div>
    </div>
</main>

<style>
    .app-item {
        transition: transform 0.2s;
    }
    .app-item:hover {
        transform: translateY(-2px);
    }
</style>
