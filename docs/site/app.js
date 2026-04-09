/* ============================================================
   CRM-Service Showcase — Interactive Functionality
   Zero dependencies. Pure vanilla JS.
   ============================================================ */

(function () {
  'use strict';

  // ---- Mock Data ----
  const MOCK_STUDENTS = [
    { id: 1, name: 'Aisha Rahman', email: 'aisha@school.edu', group: 'Web Dev A', status: 'active' },
    { id: 2, name: 'Marcus Chen', email: 'marcus@school.edu', group: 'Data Science B', status: 'active' },
    { id: 3, name: 'Sofia Petrov', email: 'sofia@school.edu', group: 'Mobile Dev A', status: 'active' },
    { id: 4, name: 'James Okafor', email: 'james@school.edu', group: 'Web Dev B', status: 'pending' },
    { id: 5, name: 'Lina Müller', email: 'lina@school.edu', group: 'UX Design A', status: 'active' },
    { id: 6, name: 'Diego Torres', email: 'diego@school.edu', group: 'Data Science A', status: 'active' },
    { id: 7, name: 'Yuki Tanaka', email: 'yuki@school.edu', group: 'Mobile Dev B', status: 'inactive' },
    { id: 8, name: 'Priya Sharma', email: 'priya@school.edu', group: 'Web Dev A', status: 'active' },
  ];

  const MOCK_COURSES = [
    { id: 1, title: 'Web Development', students: 42, fee: '$1,200', status: 'active' },
    { id: 2, title: 'Data Science', students: 35, fee: '$1,500', status: 'active' },
    { id: 3, title: 'Mobile Development', students: 28, fee: '$1,350', status: 'active' },
    { id: 4, title: 'UX Design', students: 31, fee: '$1,100', status: 'active' },
    { id: 5, title: 'Cloud Computing', students: 22, fee: '$1,400', status: 'pending' },
    { id: 6, title: 'Cybersecurity', students: 19, fee: '$1,600', status: 'active' },
  ];

  const MOCK_TEACHERS = [
    { id: 1, name: 'Dr. Emily Watson', subject: 'Web Development', groups: 3, status: 'active' },
    { id: 2, name: 'Prof. David Kim', subject: 'Data Science', groups: 2, status: 'active' },
    { id: 3, name: 'Ms. Sarah Johnson', subject: 'UX Design', groups: 2, status: 'active' },
    { id: 4, name: 'Mr. Robert Liu', subject: 'Mobile Dev', groups: 2, status: 'active' },
    { id: 5, name: 'Dr. Ana García', subject: 'Cloud Computing', groups: 1, status: 'pending' },
    { id: 6, name: 'Prof. Michael Brown', subject: 'Cybersecurity', groups: 2, status: 'active' },
  ];

  // ---- Enrollment chart data (monthly) ----
  const ENROLLMENT_DATA = [28, 35, 42, 38, 52, 61, 48, 55, 72, 68, 85, 92];
  const ENROLLMENT_LABELS = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];

  // ---- Revenue chart data ----
  const REVENUE_DATA = [4200, 5100, 4800, 6200, 7100, 6800, 7500, 8200, 7900, 8800, 9200, 10500];

  // ---- API Endpoints Data ----
  const API_ENDPOINTS = {
    'Core': [
      { method: 'GET', path: '/api/v1/teachers', desc: 'List all teachers (paginated)' },
      { method: 'POST', path: '/api/v1/teachers', desc: 'Create a new teacher' },
      { method: 'GET', path: '/api/v1/students', desc: 'List all students (paginated)' },
      { method: 'POST', path: '/api/v1/students', desc: 'Create a new student' },
      { method: 'GET', path: '/api/v1/courses', desc: 'List all courses (paginated)' },
      { method: 'POST', path: '/api/v1/courses', desc: 'Create a new course' },
      { method: 'GET', path: '/api/v1/groups', desc: 'List all groups (paginated)' },
      { method: 'PUT', path: '/api/v1/groups/:id', desc: 'Update a group' },
    ],
    'Scheduling': [
      { method: 'GET', path: '/api/v1/timetables', desc: 'List all timetables' },
      { method: 'POST', path: '/api/v1/timetables', desc: 'Create timetable entry' },
      { method: 'POST', path: '/api/v1/attendance', desc: 'Mark student attendance' },
      { method: 'GET', path: '/api/v1/attendance/group/:id', desc: 'Group attendance records' },
    ],
    'Grades': [
      { method: 'POST', path: '/api/v1/grades', desc: 'Create grade entry' },
      { method: 'GET', path: '/api/v1/grades/student/:id', desc: 'Get student grades' },
      { method: 'POST', path: '/api/v1/exams', desc: 'Create exam' },
      { method: 'GET', path: '/api/v1/exams/:id/statistics', desc: 'Get exam statistics' },
      { method: 'POST', path: '/api/v1/exams/:id/results', desc: 'Submit exam result' },
    ],
    'Finance': [
      { method: 'GET', path: '/api/v1/payments', desc: 'List all payments' },
      { method: 'POST', path: '/api/v1/payments', desc: 'Create payment record' },
      { method: 'GET', path: '/api/v1/invoices', desc: 'List all invoices' },
      { method: 'POST', path: '/api/v1/invoices', desc: 'Generate invoice' },
      { method: 'GET', path: '/api/v1/invoices/student/:id', desc: 'Student invoices' },
    ],
    'Communication': [
      { method: 'POST', path: '/api/v1/notifications', desc: 'Send notification' },
      { method: 'GET', path: '/api/v1/notifications/user/:id', desc: 'User notifications' },
      { method: 'POST', path: '/api/v1/messages', desc: 'Send message' },
      { method: 'GET', path: '/api/v1/messages/inbox', desc: 'Get inbox messages' },
    ],
    'Analytics': [
      { method: 'GET', path: '/api/v1/analytics/dashboard', desc: 'Dashboard metrics' },
      { method: 'GET', path: '/api/v1/analytics/financial', desc: 'Financial analytics' },
      { method: 'GET', path: '/api/v1/analytics/student-progress/:id', desc: 'Student progress' },
      { method: 'GET', path: '/api/v1/analytics/attendance/:id', desc: 'Attendance analytics' },
    ],
  };

  // ============================================================
  // NAVIGATION
  // ============================================================
  function initNavigation() {
    const navbar = document.getElementById('navbar');
    const mobileToggle = document.getElementById('mobile-toggle');
    const navLinks = document.getElementById('nav-links');

    // Scroll effect
    let lastScroll = 0;
    window.addEventListener('scroll', () => {
      const currentScroll = window.pageYOffset;
      if (currentScroll > 50) {
        navbar.classList.add('scrolled');
      } else {
        navbar.classList.remove('scrolled');
      }
      lastScroll = currentScroll;
    }, { passive: true });

    // Mobile toggle
    if (mobileToggle) {
      mobileToggle.addEventListener('click', () => {
        navLinks.classList.toggle('mobile-open');
        mobileToggle.textContent = navLinks.classList.contains('mobile-open') ? '✕' : '☰';
      });
    }

    // Smooth scroll for nav links
    document.querySelectorAll('a[href^="#"]').forEach(link => {
      link.addEventListener('click', (e) => {
        e.preventDefault();
        const target = document.querySelector(link.getAttribute('href'));
        if (target) {
          target.scrollIntoView({ behavior: 'smooth', block: 'start' });
          // Close mobile menu
          if (navLinks) navLinks.classList.remove('mobile-open');
          if (mobileToggle) mobileToggle.textContent = '☰';
        }
      });
    });
  }

  // ============================================================
  // SCROLL REVEAL ANIMATIONS
  // ============================================================
  function initScrollReveal() {
    const observer = new IntersectionObserver((entries) => {
      entries.forEach(entry => {
        if (entry.isIntersecting) {
          entry.target.classList.add('visible');
        }
      });
    }, {
      threshold: 0.1,
      rootMargin: '0px 0px -50px 0px'
    });

    document.querySelectorAll('.reveal').forEach(el => observer.observe(el));
  }

  // ============================================================
  // COUNTER ANIMATIONS
  // ============================================================
  function animateCounter(el, target, duration = 2000) {
    let start = 0;
    const startTime = performance.now();
    const suffix = el.dataset.suffix || '';
    const prefix = el.dataset.prefix || '';

    function update(currentTime) {
      const elapsed = currentTime - startTime;
      const progress = Math.min(elapsed / duration, 1);
      // Ease out cubic
      const eased = 1 - Math.pow(1 - progress, 3);
      const current = Math.round(start + (target - start) * eased);

      el.textContent = prefix + current.toLocaleString() + suffix;

      if (progress < 1) {
        requestAnimationFrame(update);
      }
    }

    requestAnimationFrame(update);
  }

  function initCounters() {
    const observer = new IntersectionObserver((entries) => {
      entries.forEach(entry => {
        if (entry.isIntersecting && !entry.target.dataset.animated) {
          entry.target.dataset.animated = 'true';
          const target = parseInt(entry.target.dataset.target, 10);
          animateCounter(entry.target, target);
        }
      });
    }, { threshold: 0.5 });

    document.querySelectorAll('[data-counter]').forEach(el => observer.observe(el));
  }

  // ============================================================
  // DASHBOARD INTERACTIVITY
  // ============================================================
  let currentView = 'students';

  function initDashboard() {
    // Sidebar navigation
    document.querySelectorAll('.dash-nav a[data-view]').forEach(link => {
      link.addEventListener('click', (e) => {
        e.preventDefault();
        const view = link.dataset.view;
        switchDashboardView(view);
      });
    });

    // Search functionality
    const searchInput = document.getElementById('dash-search-input');
    if (searchInput) {
      searchInput.addEventListener('input', (e) => {
        filterDashboardTable(e.target.value);
      });
    }

    // Initial render
    renderDashboard('students');
    renderChart();
  }

  function switchDashboardView(view) {
    currentView = view;

    // Update active nav
    document.querySelectorAll('.dash-nav a').forEach(a => a.classList.remove('active'));
    const activeLink = document.querySelector(`.dash-nav a[data-view="${view}"]`);
    if (activeLink) activeLink.classList.add('active');

    // Update title
    const titleMap = {
      students: 'Students',
      courses: 'Courses',
      teachers: 'Teachers'
    };
    const titleEl = document.getElementById('dash-view-title');
    if (titleEl) titleEl.textContent = titleMap[view] || 'Overview';

    // Update search placeholder
    const searchInput = document.getElementById('dash-search-input');
    if (searchInput) searchInput.placeholder = `Search ${titleMap[view] || ''}...`;
    if (searchInput) searchInput.value = '';

    // Render table
    renderDashboard(view);
  }

  function renderDashboard(view) {
    const tableHead = document.getElementById('dash-table-head');
    const tableBody = document.getElementById('dash-table-body');
    if (!tableHead || !tableBody) return;

    let headHtml = '';
    let bodyHtml = '';

    if (view === 'students') {
      headHtml = '<tr><th>Student</th><th>Email</th><th>Group</th><th>Status</th></tr>';
      bodyHtml = MOCK_STUDENTS.map(s => `
        <tr>
          <td><span class="avatar-sm">${s.name.split(' ').map(n => n[0]).join('')}</span>${s.name}</td>
          <td>${s.email}</td>
          <td>${s.group}</td>
          <td><span class="status-badge status-${s.status}">${s.status}</span></td>
        </tr>
      `).join('');
    } else if (view === 'courses') {
      headHtml = '<tr><th>Course</th><th>Students</th><th>Monthly Fee</th><th>Status</th></tr>';
      bodyHtml = MOCK_COURSES.map(c => `
        <tr>
          <td><span class="avatar-sm">📚</span>${c.title}</td>
          <td>${c.students}</td>
          <td>${c.fee}</td>
          <td><span class="status-badge status-${c.status}">${c.status}</span></td>
        </tr>
      `).join('');
    } else if (view === 'teachers') {
      headHtml = '<tr><th>Teacher</th><th>Subject</th><th>Groups</th><th>Status</th></tr>';
      bodyHtml = MOCK_TEACHERS.map(t => `
        <tr>
          <td><span class="avatar-sm">${t.name.split(' ').pop()[0]}${t.name.split(' ')[0][0]}</span>${t.name}</td>
          <td>${t.subject}</td>
          <td>${t.groups}</td>
          <td><span class="status-badge status-${t.status}">${t.status}</span></td>
        </tr>
      `).join('');
    }

    tableHead.innerHTML = headHtml;
    tableBody.innerHTML = bodyHtml;
  }

  function filterDashboardTable(query) {
    const q = query.toLowerCase().trim();
    const rows = document.querySelectorAll('#dash-table-body tr');
    rows.forEach(row => {
      const text = row.textContent.toLowerCase();
      row.style.display = text.includes(q) ? '' : 'none';
    });
  }

  // ============================================================
  // CHART RENDERING (Canvas API)
  // ============================================================
  function renderChart() {
    const canvas = document.getElementById('enrollment-chart');
    if (!canvas) return;

    const ctx = canvas.getContext('2d');
    const dpr = window.devicePixelRatio || 1;
    const rect = canvas.parentElement.getBoundingClientRect();

    canvas.width = rect.width * dpr;
    canvas.height = rect.height * dpr;
    ctx.scale(dpr, dpr);

    const W = rect.width;
    const H = rect.height;
    const padding = { top: 20, right: 20, bottom: 30, left: 40 };
    const chartW = W - padding.left - padding.right;
    const chartH = H - padding.top - padding.bottom;

    const maxVal = Math.max(...ENROLLMENT_DATA) * 1.15;
    const stepCount = ENROLLMENT_DATA.length;

    // Clear
    ctx.clearRect(0, 0, W, H);

    // Grid lines
    ctx.strokeStyle = 'rgba(255,255,255,0.05)';
    ctx.lineWidth = 1;
    for (let i = 0; i <= 4; i++) {
      const y = padding.top + (chartH / 4) * i;
      ctx.beginPath();
      ctx.moveTo(padding.left, y);
      ctx.lineTo(W - padding.right, y);
      ctx.stroke();
    }

    // Y-axis labels
    ctx.fillStyle = 'rgba(255,255,255,0.3)';
    ctx.font = '10px Inter, sans-serif';
    ctx.textAlign = 'right';
    for (let i = 0; i <= 4; i++) {
      const y = padding.top + (chartH / 4) * i;
      const val = Math.round(maxVal - (maxVal / 4) * i);
      ctx.fillText(val, padding.left - 8, y + 4);
    }

    // X-axis labels
    ctx.textAlign = 'center';
    ENROLLMENT_LABELS.forEach((label, i) => {
      const x = padding.left + (chartW / (stepCount - 1)) * i;
      ctx.fillText(label, x, H - 8);
    });

    // Build points
    const points = ENROLLMENT_DATA.map((val, i) => ({
      x: padding.left + (chartW / (stepCount - 1)) * i,
      y: padding.top + chartH - (val / maxVal) * chartH
    }));

    // Gradient fill
    const gradient = ctx.createLinearGradient(0, padding.top, 0, padding.top + chartH);
    gradient.addColorStop(0, 'rgba(59,130,246,0.25)');
    gradient.addColorStop(1, 'rgba(59,130,246,0)');

    ctx.beginPath();
    ctx.moveTo(points[0].x, padding.top + chartH);
    points.forEach((p, i) => {
      if (i === 0) {
        ctx.lineTo(p.x, p.y);
      } else {
        // Smooth curve
        const prev = points[i - 1];
        const cpX = (prev.x + p.x) / 2;
        ctx.bezierCurveTo(cpX, prev.y, cpX, p.y, p.x, p.y);
      }
    });
    ctx.lineTo(points[points.length - 1].x, padding.top + chartH);
    ctx.closePath();
    ctx.fillStyle = gradient;
    ctx.fill();

    // Line
    const lineGradient = ctx.createLinearGradient(padding.left, 0, W - padding.right, 0);
    lineGradient.addColorStop(0, '#3b82f6');
    lineGradient.addColorStop(1, '#8b5cf6');

    ctx.beginPath();
    points.forEach((p, i) => {
      if (i === 0) {
        ctx.moveTo(p.x, p.y);
      } else {
        const prev = points[i - 1];
        const cpX = (prev.x + p.x) / 2;
        ctx.bezierCurveTo(cpX, prev.y, cpX, p.y, p.x, p.y);
      }
    });
    ctx.strokeStyle = lineGradient;
    ctx.lineWidth = 2.5;
    ctx.stroke();

    // Dots
    points.forEach(p => {
      ctx.beginPath();
      ctx.arc(p.x, p.y, 3.5, 0, Math.PI * 2);
      ctx.fillStyle = '#3b82f6';
      ctx.fill();
      ctx.strokeStyle = 'rgba(59,130,246,0.3)';
      ctx.lineWidth = 4;
      ctx.stroke();
    });
  }

  // Resize handler
  let resizeTimer;
  window.addEventListener('resize', () => {
    clearTimeout(resizeTimer);
    resizeTimer = setTimeout(renderChart, 200);
  });

  // ============================================================
  // API EXPLORER
  // ============================================================
  function initApiExplorer() {
    renderApiEndpoints('Core');

    // Category buttons
    document.querySelectorAll('.api-cat-btn').forEach(btn => {
      btn.addEventListener('click', () => {
        document.querySelectorAll('.api-cat-btn').forEach(b => b.classList.remove('active'));
        btn.classList.add('active');
        renderApiEndpoints(btn.dataset.category);
      });
    });
  }

  function renderApiEndpoints(category) {
    const container = document.getElementById('api-endpoints-list');
    if (!container) return;

    const endpoints = API_ENDPOINTS[category] || [];
    container.innerHTML = endpoints.map((ep, i) => {
      const methodClass = `method-${ep.method.toLowerCase()}`;
      return `
        <div class="api-endpoint" onclick="this.classList.toggle('expanded')">
          <span class="api-method ${methodClass}">${ep.method}</span>
          <span class="api-path">${ep.path}</span>
          <span class="api-desc">${ep.desc}</span>
          <span class="api-expand-icon">▾</span>
        </div>
        <div class="api-detail">
          <div class="api-code">${generateApiExample(ep)}</div>
        </div>
      `;
    }).join('');
  }

  function generateApiExample(ep) {
    if (ep.method === 'GET') {
      return `<span class="keyword">curl</span> -H <span class="string">"X-API-Key: your-key"</span> \\
  <span class="string">"http://localhost:8080${ep.path}?page=1&page_size=10"</span>

<span class="keyword">Response:</span>
{
  <span class="string">"success"</span>: <span class="keyword">true</span>,
  <span class="string">"data"</span>: [...],
  <span class="string">"total"</span>: <span class="number">42</span>,
  <span class="string">"page"</span>: <span class="number">1</span>,
  <span class="string">"total_pages"</span>: <span class="number">5</span>
}`;
    } else if (ep.method === 'POST') {
      return `<span class="keyword">curl</span> -X POST \\
  -H <span class="string">"X-API-Key: your-key"</span> \\
  -H <span class="string">"Content-Type: application/json"</span> \\
  -d <span class="string">'{"name": "Example", "email": "ex@school.edu"}'</span> \\
  <span class="string">"http://localhost:8080${ep.path}"</span>

<span class="keyword">Response:</span>
{
  <span class="string">"success"</span>: <span class="keyword">true</span>,
  <span class="string">"message"</span>: <span class="string">"Created successfully"</span>,
  <span class="string">"data"</span>: { <span class="string">"id"</span>: <span class="string">"uuid-here"</span> }
}`;
    } else if (ep.method === 'PUT') {
      return `<span class="keyword">curl</span> -X PUT \\
  -H <span class="string">"X-API-Key: your-key"</span> \\
  -H <span class="string">"Content-Type: application/json"</span> \\
  -d <span class="string">'{"name": "Updated Name"}'</span> \\
  <span class="string">"http://localhost:8080${ep.path}"</span>

<span class="keyword">Response:</span>
{
  <span class="string">"success"</span>: <span class="keyword">true</span>,
  <span class="string">"message"</span>: <span class="string">"Updated successfully"</span>
}`;
    } else {
      return `<span class="keyword">curl</span> -X DELETE \\
  -H <span class="string">"X-API-Key: your-key"</span> \\
  <span class="string">"http://localhost:8080${ep.path}"</span>

<span class="keyword">Response:</span>
{
  <span class="string">"success"</span>: <span class="keyword">true</span>,
  <span class="string">"message"</span>: <span class="string">"Deleted successfully"</span>
}`;
    }
  }

  // ============================================================
  // INITIALIZATION
  // ============================================================
  document.addEventListener('DOMContentLoaded', () => {
    initNavigation();
    initScrollReveal();
    initCounters();
    initDashboard();
    initApiExplorer();
  });

})();
