**For Each Cycle:**

1. **Discovery Phase**
   - **Process the document sequentially using the Validation Checklist as your guide**. Address one checklist item at a time by completing all corresponding sections in the document before moving to the next item
   - **For EACH section, identify ALL activities needed** based on what information is missing or unclear. Consider relevant research areas, best practices, edge cases, and success criteria
   - **ALWAYS launch multiple specialist agents in parallel** to investigate the identified activities. Select agents based on the type of research needed (analysis, research, clarification, etc.)
   - **After receiving user feedback, identify NEW research needs** based on their input and launch additional specialist agents to investigate any new questions or directions

2. **Documentation Phase**
   - Structure findings into supporting documentation:
     - Create/update docs/domain/ for business rules and domain patterns discovered
     - Create/update docs/patterns/ for reusable solution patterns identified
     - Create/update docs/interfaces/ for external service contracts needed
   - Update the main document:
     - Base content on research findings gathered from specialist agents
     - Incorporate user feedback and additional research conducted
     - Apply command-specific document update rules
     - Focus only on current section/area being processed

3. **Review Phase**
   - **Present ALL agent findings to the user** including:
     - Complete responses from each agent (not summaries)
     - Conflicting information or recommendations
     - Proposed content based on the research
     - Questions that need user clarification
   - Show what documentation was created/updated with file paths
   - Present what was added to main document, what questions remain, and ask if you should continue
   - **Wait for user confirmation** before proceeding to next cycle

**🤔 Ask yourself each cycle:**
1. **Discovery**: Have I identified ALL activities needed for this section/area?
2. **Discovery**: Have I launched parallel specialist agents to investigate?
3. **Documentation**: Have I created/updated supporting documentation (domain/patterns/interfaces)?
4. **Documentation**: Have I updated the main document according to command-specific rules?
5. **Review**: Have I presented COMPLETE agent responses to the user (not summaries)?
6. **Review**: Have I received user confirmation before proceeding to next cycle?
7. Are there more sections/areas that need investigation?
8. Should I continue to the next section/area or wait for user input?
9. If work is complete, have I asked the user for confirmation to proceed?